import { useState } from 'react';
import { useRef } from 'react';
import axios from "axios";

function Login(){
  const emailRef = useRef(null);
  const passwordRef = useRef(null);
  const [email, setEmail] = useState("  ");
  const [password, setPassword] = useState("");
  const [errorMsg, setErrorMsg] = useState("");

  var emailHandler = (event) => {
    setEmail(event.target.value);
  }

  var passwordHandler = (event) => {
    setPassword(event.target.value);
  }

  var loginHandler = () => {
    axios.post(
      "http://"+document.location.host+"/api/authorize",
      {
        "useremail": email,
        "password": password
      }
    )
    .then((response) => {
      if (response.status == 200){
        setErrorMsg("")
        emailRef?.current?.classList.remove("is-invalid")
        passwordRef?.current?.classList.remove("is-invalid")
        window.location.href = "/home"
      }
    })
    .catch((error) => {
      if (error.response.status == 400){
        setErrorMsg("An Account with this email does not exist!")
        emailRef?.current?.classList.add("is-invalid")
        passwordRef?.current?.classList.remove("is-invalid")
      }else if (error.response.status == 401){
        setErrorMsg("Wrong password!")
        emailRef?.current?.classList.remove("is-invalid")
        passwordRef?.current?.classList.add("is-invalid")
      }else {
        emailRef?.current?.classList.remove("is-invalid")
        passwordRef?.current?.classList.remove("is-invalid")
        setErrorMsg("Server error!")
      }
    });
  }

  return (
      <main class="form-signin w-100 m-auto">
      <div class="form-floating">
          <input ref={emailRef} type="email" class="form-control mb-2 rounded" id="floatingInput" placeholder="name@example.com" onChange={emailHandler}/>
          <label for="floatingInput">Email</label>
      </div>
      <div class="form-floating">
          <input ref={passwordRef} type="password" class="form-control rounded" id="floatingPassword" placeholder="password" onChange={passwordHandler}/>
          <label for="floatingPassword">Password</label>
      </div>
      <button class="btn btn-primary w-100 py-2 mb-2" onClick={loginHandler}>Sign in</button>
      <a class="btn btn-primary w-100 py-2 mb-2" href="/register" >Register</a>
      { 
        errorMsg 
        ? <div class="alert alert-danger w-100 py-2" role="alert">{errorMsg}</div> 
        : <></>
      }
      </main>
  );
}

function Register(){
  return (
      <main class="form-signin w-100 m-auto">
      <div id="regForm">
          <div class="form-floating">
          <input type="email" class="form-control mb-2 rounded" id="floatingInput" placeholder="name@example.com"/>
          <label for="floatingInput">Email</label>
          </div>
          <div class="form-floating">
          <input type="username" class="form-control mb-2 rounded" id="floatingInput" placeholder="Username"/>
          <label for="floatingInput">Username</label>
          </div>
          <div class="form-floating">
          <input type="password" class="form-control mb-2 rounded" id="floatingPassword" placeholder="Password"/>
          <label for="floatingPassword">Password</label>
          </div>
          <div class="form-floating">
          <input type="password" class="form-control rounded" id="floatingPassword" placeholder="Password"/>
          <label for="floatingPassword">Confirm Password</label>
          </div>
          <button class="btn btn-primary w-100 py-2 mb-2">Register</button>
      </div>
  </main>
  );
}

function App() {
  var appState = "login"

  switch(appState){
    case "login":
      return (
        <div className="App">
          <Login/>
        </div>
      );
    case "register":
      return (
        <div className="App">
          <Register/>
        </div>
      );
    default:
      return
  }
}
  
export default App;