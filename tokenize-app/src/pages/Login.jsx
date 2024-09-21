import React, { useState } from "react";
import { replace, useNavigate } from "react-router-dom";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";
// import axios from 'axios'
import { loginApiCall } from "../services/ApiCaller";

export const Login = () => {
  const successNotify = (msg) => toast.success(msg, { icon: true });
  const errorNotify = (err) => toast.error(err, { icon: true });

  const navigate = useNavigate();

  //state variables
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");

  //to show loading and success states
  const [isLoading, setLoading] = useState(false);
  const [message, setMessage] = useState("");

  //handlers
  const handleNavigateToSignUp = () => {
    navigate("/signup");
  };
  const handleLoginEmail = (e) => {
    setEmail(e.target.value);
  };
  const handleLoginPass = (e) => {
    setPassword(e.target.value);
  };
  const handleSumbit = async (e) => {
    setTimeout(()=>{
      console.log("delay")
    },1500);
    e.preventDefault();

    if (email === "") {
      setMessage("error");
      errorNotify("Email is required !");
      return;
    }
    if (password === "") {
      setMessage("error");

      errorNotify("Password is required!");
      return;
    } else {
      console.log("comthi");
    }

    setLoading(true);

    const response = await loginApiCall(email, password);
    if (response === true) {
      console.log("Login success !");
      setMessage("success");
      successNotify(email + " has logged in successfully 🚀");
      navigate('/dashboard',true);
    } else {
      console.log("No response font-end");
      setMessage("error");
      errorNotify(email + "either exists or need to register");
    }

    setLoading(false);
  };

  return (
    <div className="flex  min-h-full flex-col justify-center px-6 py-12 lg:px-8">
      <img
        className="mx-auto h-10 w-auto"
        src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600"
      />

      <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
        Login Page
      </h2>

      {/* form  */}
      <div>
        <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
          <form
            className="space-y-6"
            onSubmit={handleSumbit}
            method="POST" //change it
          >
            <div>
              <label
                htmlFor="email"
                className="block text-sm font-medium leading-6 text-gray-900"
              >
                Email Address
              </label>

              <div className="mt-2">
                <input
                  id="email-id"
                  name="Email"
                  type="email"
                  onChange={handleLoginEmail}
                  value={email}
                  required
                  placeholder="abc@gmail.com"
                  className="block w-full rounded-md border-0 py-1.5 placeholder:px-2 text-gray-900 shadow-sm ring-1 ring-gray-300 placeholder:text-gray-400 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>
            </div>

            {/* password  */}
            <div>
              <div className="flex items-center justify-between">
                <label
                  htmlFor="password"
                  className="block text-sm font-medium leading-6 text-gray-900"
                >
                  Password
                </label>
              </div>

              <div className="mt-2">
                <input
                  id="password-id"
                  name="Password"
                  value={password}
                  required
                  type="password"
                  onChange={handleLoginPass}
                  placeholder="****3**"
                  className="block w-full rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 px-4 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                />
              </div>
            </div>

            <div>
              <button
                type="submit"
                className="flex w-full justify-center rounded-md bg-indigo-600 leading-10 px-3 py-1.5 text-sm font-semibold text-white shadow-sm hover: bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                disabled={isLoading ? true : false}
              >
                {isLoading ? "Loading..." : "Login"}
              </button>
            </div>
          </form>

          <p className="mt-10 text-center text-sm text-gray-500">
            Don't have an account ?{" "}
            <a
              href="#"
              onClick={handleNavigateToSignUp}
              className="font-semibold px-4 leading-6 text-indigo-600 hover:text-indigo-500"
            >
              Sign Up
            </a>
          </p>
        </div>
      </div>
      {message && <ToastContainer theme="dark" autoClose={3000}/>}
    </div>
  );
};
