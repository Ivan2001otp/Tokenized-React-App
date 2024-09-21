import axios from 'axios'
import { MAIN_URL,LOGIN_END_POINT } from '../Constants/config';

export const loginApiCall = async(email,password)=>{
    try{
        console.log("api-call-section")
        console.log(MAIN_URL+LOGIN_END_POINT)
        const response = await axios.post(MAIN_URL+LOGIN_END_POINT,{email,password});
        console.log(response.data);
        console.log(response.status);
        return response;
    }catch(error){
        console.log(response.data);
        console.log(response.status);
        console.error(error);
        throw error;
    }
}