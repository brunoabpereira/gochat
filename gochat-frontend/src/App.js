import { useState } from 'react';
import { useRef } from 'react';
import { useEffect } from 'react';

function Input({innerRef}){
  return (
    <input ref={innerRef} class="px-0 form-control" style={{ width: 'max-content' }}></input>
  );
}

function SendButton({sendHandler}){
  return (
    <button class="btn btn-primary py-2 mb-2" onClick={sendHandler}>send</button>
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
      <div>
          {msg.Userid}
        </div>
        <div>
          {msg.Messagetime}
        </div>
      </div>
      <div class="card-body">
        <p>{msg.Messagetext}</p>
      </div>
    </div>
  );
}

function Channel({channelName}){
  return (
    <h3 class="bg-body rounded py-1">{channelName}</h3>
  );
}

function Chat(){
  const inputRef = useRef(null);
  const [msgs, setMsgs] = useState([]);

  const ws = new WebSocket('ws://localhost:9000');

  const join = {
    Op: "join",
    Value: "1",
  };

  ws.onopen = () => {
    ws.send(JSON.stringify(join));
  };
  
  ws.onmessage = (event) => {
    const response = JSON.parse(event.data);
    if (response) {
      setMsgs(oldMsgs => [...oldMsgs, response]);
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

  return (
    <div class="container">
      <div class="row">
        <Channel channelName={"channel#4"} /> 
      </div>
      <div class="row">
        <ChatBox messagesArray={msgs}/> 
      </div>
      <div class="row">
        <div class="col">
          <Input innerRef={inputRef}/>
        </div>
        <div class="col">
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
