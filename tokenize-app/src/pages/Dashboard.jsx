import React from 'react'
import {logoutApiCall,deleteUserApiCall} from '../services/ApiCaller'

export const Dashboard = () => {

    const handleLogout=async()=>{
        console.log("logout pressed");
        let logOutResult = await logoutApiCall();
        if(logOutResult){
          console.log("User logged out");
        }else{
          console.log("logout result is false");
        }
    }

    const handleDeleteUser=async()=>{

      console.log("Invoked handleDeleteUser");
      let deleteResult = await deleteUserApiCall();
      if (deleteResult){
        console.log("User deleted!");
      }else{
        console.log("User could not be deleted!");
      }
    }

  return (
    <div
        className='flex min-h-full flex-col px-6 py-12 lg:px-8 justify-center'
    >
         <div className="sm:mx-auto sm:w-full sm:max-w-sm">
        <img
          className="mx-auto h-10 w-auto"
          src="https://tailwindui.com/img/logos/mark.svg?color=indigo&shade=600"
        />

        <h2 className="mt-10 text-center text-2xl font-bold leading-9 tracking-tight text-gray-900">
         Dashboard page
        </h2>
      </div>
      <div
      className='sm:mx-auto sm:w-full sm:max-w-sm'
      >
            <div
                className='flex justify-between w-full mt-20 items-baseline'
            >
                <button
                    onClick={handleLogout}
                    className='rounded-md bg-indigo-500 leading-10 text-sm font-semibold text-white shadow-sm hover:bg-indigo-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-indigo-500 px-8'
                >Logout</button>
                
                <button
                onClick={handleDeleteUser}
                    className='rounded-md bg-red-600 leading-10 text-sm font-semibold text-white shadow-sm hover:bg-red-600 focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-red-700 px-4'
                >Delete User</button>
            </div>
      </div>
    </div>
  )
}
