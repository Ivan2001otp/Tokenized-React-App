import React from "react";
import { useNavigate } from "react-router-dom";
import { useState } from "react";
import {  signUpApiCall ,navigateToDashboardApiCall} from "../services/ApiCaller";
import { ToastContainer, toast } from "react-toastify";
import "react-toastify/dist/ReactToastify.css";

export const SignUp = () => {
  const navigate = useNavigate();

  //notification
  const successNotify = (msg) => toast.success(msg, { icon: true });
  const errorNotify = (err) => toast.error(err, { icon: true });

  //to show loading and success states
  const [isLoading, setLoading] = useState(false);
  const [message, setMessage] = useState("");


  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [first_name, setFirstName] = useState("");
  const [last_name, setLastName] = useState("");
  const [user_name, setUserName] = useState("");
  const [phone, setPhone] = useState("");
  const [user_type, setUserType] = useState("USER");

  const handleNavigateToLogin = () => {
    navigate("/login");
  };

  const handleUserType = (e) => {
    setUserType(e.target.value);
  };

  const submitSignUpDetails = async(e) => {
    e.preventDefault();

    console.log(
      email +
        " " +
        password +
        " " +
        first_name +
        " " +
        last_name +
        " " +
        user_name +
        " " +
        phone +
        " " +
        user_type
    );

    if(email===""){
        setMessage("error");
        errorNotify("Email is required!")
    }if(password===""){
        setMessage("error");
        errorNotify("Password is required!")

    }if(first_name===""){
        setMessage("error");
        errorNotify("First-Name is required!")

    }if(last_name===""){
        setMessage("error");
        errorNotify("Last-Name is required!")

    }if(user_name===""){
        setMessage("error");
        errorNotify("UserName is required!")

    }if(phone===""){
        setMessage("error");
        errorNotify("Phone is required!")

    }if(user_type===""){
        setMessage("error");
        errorNotify("User-Type is required!")

    }

    setLoading(true);

    const response = await signUpApiCall(user_name,first_name,last_name,email,password,phone,user_type);

    if(response===true){
        console.log("success sign up");
        setMessage("success");
        successNotify("Signed Up successfully 🚀");
        //make an api call to invoke the dashboard.
      const navigateResponse = await navigateToDashboardApiCall();

      if(navigateResponse){
        navigate('/dashboard',true);
        console.log("success navigation to dashboard from login !")
      }else{
        console.log("failure navigation to dashboard from login !")
      }
    }else{
        console.log("failure sign up");
        setMessage("error");
        errorNotify("Something went wrong");
    }

    setLoading(false);

  };

  return (
    <div className="flex min-h-full flex-col justify-center px-6 py-12 lg:px-8">
      <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <img
          className="mx-auto h-10 w-auto"
          src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600"
        />

        <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
          Sign Up Page
        </h2>
      </div>

      {/* form  */}
      <div className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm">
        <form
          className="space-y-6"
          onSubmit={submitSignUpDetails}
          method="POST" //<-change it
        >
          {/* username  */}
          <div>
            <label
              htmlFor="User Name"
              className="block text-sm font-medium leading-6 text-gray-900 "
            >
              UserName
            </label>
            
            
            <input
              name="User Name"
              type="text"
              placeholder="Blaxk Coker"
              id="user-name"
              value={user_name}
              required
              className="block w-full 
                                
                                rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 
                                px-4
                                focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            
              onChange={(e) => {
                setUserName(e.target.value);
              }}
            />
          </div>

          {/* first name  */}
          <div>
            <label
              htmlFor="First Name"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              First Name
            </label>
           <input
             placeholder="David.."
            name="First Name"
            type="text"
            minLength={3}
            maxLength={20}
            className="block w-full 
                                
                                rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 
                                px-4
                                focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            
            value={first_name}
            onChange={(e)=>{setFirstName(e.target.value);}}
            id="first-name"
           />
          </div>

          {/* last name  */}
          <div>
            <label
              htmlFor="Last Name"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              Last Name
            </label>
              <input
                autoComplete="Last Name"
                id="last-name"
                type="text"
                name="Last Name"
                minLength={3}
                required
                placeholder="Beckmar.."
                maxLength={20}
                value={last_name}
                className="block w-full 
                                
                                rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 
                                px-4
                                focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            
                onChange={(e)=>{setLastName(e.target.value);}}
              />
          </div>

          {/* email  */}
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
                autoComplete="email"
                required
                placeholder="abc@gmail.com"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                className="block w-full rounded-md border-0 py-1.5 
                            placeholder:px-2
                            text-gray-900 shadow-sm ring-1 ring-gray-300 placeholder:text-gray-400 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
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
                id="password"
                name="password"
                type="password"
                required
                autoComplete="current-password"
                placeholder="***3**"
                value={password}
                onChange={(e) => setPassword(e.target.value)}
                className="block w-full 
                                
                                rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 
                                px-4
                                focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
              />
            </div>
          </div>

          {/* phone  */}
          <div>
            <label
              htmlFor="Phone"
              className="block text-sm font-medium leading-6 text-gray-900"
            >
              Phone
            </label>
            <input
                id="phone-id"
                value={phone}
                type="number"
                onChange={(e)=>{

                        setPhone(e.target.value);
                    
                }}
                required
                maxLength={10}
                minLength={10}
                placeholder="MobileNumber"
                className="block w-full 
                                 rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 
                                px-4
                                focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            
            />
          </div>

          {/* user type  dropdown*/}
          <div>
            <label
              htmlFor="User-Type"
              className="block pb-1 text-sm font-medium leading-6 text-gray-900"
            >
              User-Type
            </label>
            <select
                id="user-type"
              value={user_type}
              className="block w-full  rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 
                                px-4
                                focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
            
              onChange={handleUserType}
            >
              <option value="ADMIN">ADMIN</option>
              <option value="USER">USER</option>
            </select>
          </div>

          {/* //button */}
          <div>
            <button
              type="submit"
              disabled={isLoading ? true : false}
              className="flex w-full justify-center rounded-md bg-indigo-600 leading-10 px-3 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
            >
             {isLoading ? 'Loading...' : 'Sign Up'}
            </button>
          </div>
        </form>

        <p className="mt-10 text-center text-sm text-gray-500">
          Already have an account ?{"  "}
          <a
            href="#"
            onClick={handleNavigateToLogin}
            className="font-semibold px-4 leading-6 text-indigo-600 hover:text-indigo-500"
          >
            Login
          </a>
        </p>
      </div>
      {message && <ToastContainer theme="dark" autoClose={3000}/>}
    </div>
  );
};
