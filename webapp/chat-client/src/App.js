import { useState } from 'react';
import { useRef } from 'react';
import { useEffect } from 'react';
import axios from "axios";

function UserDropDown({appState}){
  return (
    <div class="dropdown">
      <a href="#" class="d-flex align-items-center text-white text-decoration-none dropdown-toggle" data-bs-toggle="dropdown" aria-expanded="false">
        <img src='/static/assets/person-circle.svg' width="20" height="20"></img>
        <strong class="p-2" >{appState.username}</strong>
      </a>
      <ul class="dropdown-menu dropdown-menu-dark text-small shadow">
        <li class=""><p class="mx-1 px-2 bg-body rounded">{appState.useremail}</p></li>
        <li><a class="dropdown-item" href="#">Settings</a></li>
        <li><a class="dropdown-item" href="#">Profile</a></li>
        <li><hr class="dropdown-divider"/></li>
        <li><a class="dropdown-item" href="/logout">Sign out</a></li>
      </ul>
    </div>
  );
}

function Input({innerRef, enterHandler}){
  return (
    <input ref={innerRef} class="form-control" onKeyDown={enterHandler}></input>
  );
}

function SendButton({sendHandler}){
  return (
    <button class="btn btn-primary py-2 mb-2 float-end" onClick={sendHandler}>send</button>
  );
}

function ChatBox({ messagesArray}){
  const chatBoxRef = useRef(null);
  useEffect(() => {    
    chatBoxRef?.current?.lastElementChild?.scrollIntoView();
  });
  return (
    <div ref={chatBoxRef} class="overflow-auto mb-3 p-2 pb-0 bg-body shadow-sm rounded" style={{ maxHeight: '400px', height: '400px'}}>
      {
        messagesArray.map((msg) => <Message key={msg.Messagetime} msg={msg} />)
      }
    </div>
  );
}

function Message({ msg }){
  return (
    <div class="card mb-2">
      <div class="card-header">
        <div class="container p-0">
          <div class="row">
            <div class="col p-0">
              {msg.Username}
            </div>
            <div class="col">
              <p class="text-end mb-0">
                {msg.Timestamp}
              </p>
            </div>
          </div>
        </div>
      </div>
      <div class="card-body p-1">
        <p>{msg.Text}</p>
      </div>
    </div>
  );
}

var ws = null;

function Chat({appState, setAppState}){
  const inputRef = useRef(null);
  const [msgs, setMsgs] = useState([]);

  useEffect(() => {
    if (!ws) {
      ws = new WebSocket('ws://localhost:9000/ws')  
    }

    ws.onopen = () => {
      ws.send(JSON.stringify(joinChannel));
      ws.send(JSON.stringify(getMsgs));
    };
    
    ws.onmessage = (event) => {
      const response = JSON.parse(event.data);
      if (response) {
        if (Array.isArray(response)){
          setMsgs(oldMsgs => [...oldMsgs, ...(response.reverse())]);
        }else {
          setMsgs(oldMsgs => [...oldMsgs, response]);
        }
      }
    };
  });

  const joinChannel = {
    Op: "join",
    Value: appState.currChannel.channelid.toString(),
  };

  const getMsgs = {
    Op: "get",
    Value: "20",
  };

  var sendHandler = () => {
    var text = inputRef.current.value;
    if (text) {
      ws.send(
        JSON.stringify(
          {
          op: "send",
          value: text
          }
        )
      );
      inputRef.current.value = "";
    }
  }

  var enterHandler = (event) => {
    if (event.key !== "Enter") {
      return
    }
    var text = inputRef.current.value;
    if (text) {
      ws.send(
        JSON.stringify(
          {
          op: "send",
          value: text
          }
        )
      );
      inputRef.current.value = "";
    }
  }

  const leaveChannel = {
    Op: "leave",
    Value: appState.currChannel.channelid.toString(),
  };

  var leaveHandler = () => {
    ws.send(
      JSON.stringify(leaveChannel)
    );
    ws.close()
    ws = null
    setAppState({
      ...appState,
      view: "channels"
    })
  }

  return (
    <div class="container">
      <div class="row">
        <div class="d-flex px-0">
          <div class="py-1 pe-1 flex-grow-1">
            <h3 class="bg-body rounded py-1 px-1">{appState.currChannel.channelname}</h3>
          </div>
          <div class="py-1 ps-1 ">
            <UserDropDown appState={appState}/>
          </div>
        </div>
      </div>
      <div class="row">
        <ChatBox messagesArray={msgs}/> 
      </div>
      <div class="row">
        <div class="col">
          <Input innerRef={inputRef} enterHandler={enterHandler}/>
        </div>
        <div class="col col-sm-2">
          <SendButton sendHandler={sendHandler}/>
        </div>
        <div class="col col-sm-2">
          <button class="btn btn-secondary py-2 mb-2 float-end" onClick={leaveHandler}>leave</button>
        </div>
      </div>
    </div>
  );
}

function AddChannel(){
  return (
    <div class="card text-body-secondary mt-2">
    <div class="container p-1">
      <div class="row">
        <div class="col">
          <button type="button" class="btn p-0 m-2">
          <img src='/static/assets/plus-square-dotted.svg' width="50" height="50"></img>
          </button>
        </div>
      </div>
    </div>
    </div>
  );
}

function ChannelItem({appState, setAppState, channel}){
  var channelHandler = () => {
    setAppState({
      ...appState,
      view: "chat",
      currChannel: {...channel}
    })
  }
  return (
    <div class="card text-body-secondary mt-2">
      <div class="container p-1">
        <div class="row">
          <div class="col col-sm-2">
            <button type="button" class="btn p-0 m-2" onClick={channelHandler}>
            <img src='/static/assets/arrow-right-square-fill.svg' width="50" height="50"></img>
            </button>
          </div>
          <div class="col my-2 pl-0">
           <p class="text-gray-dark fs-4 m-2">{channel.channelname}</p>
          </div>
        </div>
      </div>
    </div>
  );
}

function ChannelList({appState, setAppState}){
  useEffect(() => {    
    axios.get("http://localhost:8000/api/users").then((response) => {
      if (response.status == 200){
        setAppState({
          ...appState,
          username: response.data.username,
          useremail: response.data.useremail,
          channelList: response.data.channels
        })
      }
    });
  },[]);

  return (
    <main class="container-sm">
      <div class="row">
        <div class="d-flex px-0">
          <div class="py-1 pe-1 flex-grow-1">
            <h3 class="bg-body rounded py-1 px-1">My Channels</h3>
          </div>
          <div class="py-1 ps-1 ">
            <UserDropDown appState={appState}/>
          </div>
        </div>
      </div>
      <div class="row">
        <div class="overflow-auto p-3 bg-body rounded shadow-sm" style={{ maxHeight: "300px" }}>
          <AddChannel />
          {
            appState.channelList?.map((channel) => <ChannelItem key={channel.channelid} channel={channel} appState={appState} setAppState={setAppState}/>)
          }
        </div>
      </div>
    </main>
  );
}

function App() {
  const [appState, setAppState] = useState({
    view: "channels",
    username: "",
    useremail: "",
    channelList: [],
    currChannel: {},
  });

  switch(appState.view){
    case "chat":
      return (
        <div className="App">
          <Chat appState={appState} setAppState={setAppState}/>
        </div>
      );
    case "channels":
      return (
        <div className="App">
          <ChannelList  appState={appState} setAppState={setAppState}/>
        </div>
      );
    default:
      return
  }
}


export default App;
