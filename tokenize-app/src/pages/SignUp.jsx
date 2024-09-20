import React from "react";
import { useNavigate } from "react-router-dom";
import {useState} from 'react';
import {TextFieldContainer} from '../components/TextFieldContainer'

export const SignUp = () => {

    const navigate = useNavigate();

    const [email,setEmail] = useState('');
    const [password,setPassword]=useState('');
    const [first_name,setFirstName]=useState('');
    const [last_name,setLastName]=useState('');
    const [user_name,setUserName]=useState('');
    const [phone,setPhone]=useState('');
    const [user_type,setUserType]=useState('USER');

    const handleNavigateToLogin=()=>{
        navigate("/login");
    }

    const handleUserType=(e)=>{
        setUserType(e.target.value)
    }

    const submitLoginDetails=(e)=>{
        e.preventDefault();
        console.log("Email : ",email)
        console.log("Password : ",password)
        alert(email+" "+password);
        
    }

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
      <div
        className="mt-10 sm:mx-auto sm:w-full sm:max-w-sm"
      >
            <form
            className="space-y-6"
            action="#"
            method="" //<-change it 
            >       

                     {/* username  */}
                     <div>
                        <label htmlFor="User Name" className="block text-sm font-medium leading-6 text-gray-900 ">UserName</label>
                        <TextFieldContainer 
                                name="User Name"
                                type="name"
                                id="User Name"
                                placeholder={"Coder mamu."}
                                value={user_name}
                                onChange={(e)=>{
                                    setUserName(e.target.value);
                                }}
                               
                            />
                    </div>

                    {/* first name  */}
                    <div>
                            <label htmlFor="First Name" className="block text-sm font-medium leading-6 text-gray-900">First Name</label>
                            <TextFieldContainer 
                                name="First Name"
                                type="name"
                                id="First Name"
                                placeholder={"David.."}
                                value={first_name}
                                onChange={(e)=>{
                                    setFirstName(e.target.value);
                                }}
                               
                            />
                        
                    </div>

                    {/* last name  */}
                    <div>
                        <label htmlFor="Last Name" className="block text-sm font-medium leading-6 text-gray-900">Last Name</label>
                        <TextFieldContainer
                            name="Last Name"
                            type="name"
                            id="Last Name"
                            placeholder={"Feng..."}
                            value={last_name}
                            onChange={(e)=>setLastName(e.target.value)}
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
                            <input id="email" name="email" type="email" autoComplete="email" required placeholder="abc@gmail.com"
                            value={email}
                            onChange={(e)=>setEmail(e.target.value)}
                            className="block w-full rounded-md border-0 py-1.5 
                            placeholder:px-2
                            text-gray-900 shadow-sm ring-1 ring-gray-300 placeholder:text-gray-400 focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                            />
                        </div>
                    </div>

                    {/* password  */}
                    <div>
                        <div className="flex items-center justify-between">
                            <label htmlFor="password" className="block text-sm font-medium leading-6 text-gray-900">Password</label>
                            
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
                                onChange={
                                    (e)=>setPassword(e.target.value)
                                }
                                className="block w-full 
                                
                                rounded-md border-0 py-1.5 text-gray-900 shadow-sm ring-1 ring-inset ring-gray-300 placeholder:text-gray-400 focus:ring-2 
                                px-4
                                focus:ring-inset focus:ring-indigo-600 sm:text-sm sm:leading-6"
                            />
                        </div>
                    </div>

                    {/* phone  */}
                    <div>
                        <label htmlFor="Phone" className="block text-sm font-medium leading-6 text-gray-900">Phone</label>
                        <TextFieldContainer 
                                name="User Name"
                                type="name"
                                id="User Name"
                                placeholder={"8829012245"}
                                value={phone}
                                onChange={(e)=>{
                                    setPhone(e.target.value);
                                }}
                               
                            />
                    </div>

                    {/* user type  dropdown*/}
                    <div>
                    <label htmlFor="User-Type" className="block pb-1 text-sm font-medium leading-6 text-gray-900">User-Type</label>
                        <select
                            className="px-4 py-2"
                            value={user_type}
                            onChange={handleUserType}
                        >
                            <option value="ADMIN">ADMIN</option>
                            <option value="USER">USER</option>
                        </select>
                    </div>

                    {/* //button */}
                     <div>
                        <button
                        onClick={
                            submitLoginDetails
                        }
                        type="submit"
                        className="flex w-full justify-center rounded-md bg-indigo-600 leading-10 px-3 py-1.5 text-sm font-semibold text-white shadow-sm hover:bg-indigo-500 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-600"
                        >Sign Up</button>
                    </div>
            </form>

            <p className="mt-10 text-center text-sm text-gray-500">
                       Already have an account ?{'  '}
                        <a href="#" 
                        onClick={
                            handleNavigateToLogin
                        }
                        className="font-semibold px-4 leading-6 text-indigo-600 hover:text-indigo-500">
                            Login
                        </a>
            </p>
      </div>

    </div>
  );
};
