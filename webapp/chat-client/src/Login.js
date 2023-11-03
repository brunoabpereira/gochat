import { useState } from 'react';
import { useRef } from 'react';
import axios from "axios";

function Login({setAppState}){
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
    // todo: check if any field is empty
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

  var registerHandler = () => {
    setAppState("register")
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
      <button class="btn btn-primary w-100 py-2 mb-2" onClick={registerHandler}>Register</button>
      { 
        errorMsg 
        ? <div class="alert alert-danger w-100 py-2" role="alert">{errorMsg}</div> 
        : <></>
      }
      </main>
  );
}

function RegisterSucccess(){
  return (
    <div class="form-signin w-100 m-auto">
      <h5 class="bs-primary-text-success" style={{ textAlign: "center" }}>
        Success! Your account has been created.
      </h5>
      <a class="btn btn-success w-100 py-2 mb-2" href="/home">Continue</a>
    </div>
  );
}

function Register(){
  const [errorMsg, setErrorMsg] = useState("");
  const [regState, setRegState] = useState("register");
  const emailRef = useRef(null);
  const usernameRef = useRef(null);
  const passwordRef = useRef(null);
  const password2Ref = useRef(null);
  const registerRef = useRef(null);
  
  var formChange = () => {
    if (
      !emailRef?.current.value     || 
      !usernameRef?.current.value  || 
      !passwordRef?.current.value  || 
      !password2Ref?.current.value 
    ){
      registerRef?.current?.classList.add("disabled")
    }else {
      registerRef?.current?.classList.remove("disabled")

      if ( passwordRef?.current.value && password2Ref?.current.value ){
        if ( passwordRef?.current.value === password2Ref?.current.value ){
          password2Ref?.current.classList.add("is-valid")
          password2Ref?.current.classList.remove("is-invalid")
        }else {
          password2Ref?.current.classList.remove("is-valid")
          password2Ref?.current.classList.add("is-invalid")
        }
      }

      emailRef?.current?.classList.remove("is-invalid")
      usernameRef?.current?.classList.remove("is-invalid")
    }
  }

  var registerHandler = () => {
    axios.post(
      "http://"+document.location.host+"/api/users",
      {
        "useremail": emailRef?.current.value,
        "username": usernameRef?.current.value,
        "password": passwordRef?.current.value
      }
    )
    .then((response) => {
      if (response.status == 200){
          axios.post(
            "http://"+document.location.host+"/api/authorize",
            {
              "useremail": emailRef?.current.value,
              "password": passwordRef?.current.value
            }
          );
          setRegState("success")
      }
    })
    .catch((error) => {
      if (error.response?.status == 400){
        if (error.response.data["error"]){
          setErrorMsg("")
          switch(error.response.data["error"]){
            case "Username already used.":
              emailRef?.current?.classList.remove("is-invalid")
              usernameRef?.current?.classList.add("is-invalid")
              break;
            case "Email already used.":
              emailRef?.current?.classList.add("is-invalid")
              usernameRef?.current?.classList.remove("is-invalid")
              break;
          }
        }else {
          setErrorMsg("Server error!")
        }
      }
    });
  }


  switch(regState){
    case "register":
      return (
        <main class="form-signin w-100 m-auto">
        <div id="regForm" onChange={formChange}>
            <div class="form-floating">
              <input ref={usernameRef} type="username" class="form-control mb-2 rounded" id="floatingInput" placeholder="Username"/>
              <label for="floatingInput">Username</label>
              <div class="invalid-tooltip">Username already in use.</div>
            </div>
            <div class="form-floating">
              <input ref={emailRef} type="email" class="form-control mb-2 rounded" id="floatingInput" placeholder="name@example.com"/>
              <label for="floatingInput">Email</label>
              <div class="invalid-tooltip">Email already in use.</div>
            </div>
            <div class="form-floating">
              <input ref={passwordRef} type="password" class="form-control mb-2 rounded" id="floatingPassword" placeholder="Password"/>
              <label for="floatingPassword">Password</label>
            </div>
            <div class="form-floating">
              <input ref={password2Ref} type="password" class="form-control rounded" id="floatingPassword" placeholder="Password"/>
              <div class="invalid-tooltip">Passwords do not match.</div>
              <label for="floatingPassword">Confirm Password</label>
            </div>
            <button ref={registerRef} class="btn btn-primary w-100 py-2 mb-2 disabled" onClick={registerHandler}>Register</button>
            { 
              errorMsg 
              ? <div class="alert alert-danger w-100 py-2" role="alert">{errorMsg}</div> 
              : <></>
            }
        </div>
      </main>
      );
    case "success":
      return (
        <RegisterSucccess />
      );
  }
}

function App() {
  const [appState, setAppState] = useState("login");

  switch(appState){
    case "login":
      return (
        <div className="App">
          <Login setAppState={setAppState} />
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