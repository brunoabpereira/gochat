import { useState } from 'react';
import { useRef } from 'react';
import { useEffect } from 'react';

function Input({innerRef, enterHandler}){
  return (
    <input ref={innerRef} class="px-0 form-control" onKeyDown={enterHandler}></input>
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

const ws = new WebSocket('ws://localhost:9000/ws');

function Chat(){
  const inputRef = useRef(null);
  const [msgs, setMsgs] = useState([]);

  const joinChannel = {
    Op: "join",
    Value: "1",
  };

  const getMsgs = {
    Op: "get",
    Value: "20",
  };

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

  ws.onclose = () => {
    ws.close();
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
    if (event.key != "Enter") {
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
    Value: "1",
  };

  var leaveHandler = () => {
    ws.send(
      JSON.stringify(leaveChannel)
    );
    window.location.replace("/channels")
  }

  var channelName = "channel#1"
  return (
    <div class="container">
      <div class="row">
        <div class="col px-0">
          <h3 class="bg-body rounded py-1 px-1">{channelName}</h3>
        </div>
        <div class="col col-sm-2 px-0">
          <button class="btn btn-secondary py-2 mb-2 float-end" onClick={leaveHandler}>leave</button>
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
      </div>
    </div>
  );
}

function App() {
  return (
    <div className="App">
      <Chat />
    </div>
  );
}


export default App;
