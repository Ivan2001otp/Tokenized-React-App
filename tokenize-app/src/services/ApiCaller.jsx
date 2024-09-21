import axios from "axios";
import {
  MAIN_URL,
  LOGIN_END_POINT,
  REGISTER_END_POINT,
} from "../Constants/config";

    export const loginApiCall = async (email, password) => {
    try {
        console.log("api-call-section");
        console.log(MAIN_URL + LOGIN_END_POINT);
        let response = await axios.post(MAIN_URL + LOGIN_END_POINT, {
        email,
        password,
        });
        console.log("success-login");
        console.log(response);
        console.log(response.status);
        if (response.status === 200 || response.status === 201) {
        return true;
        }
        return false;
    } catch (error) {
        console.log("faileure login");
        console.error(error);
        return false;
    }
    };




export const  signUpApiCall = async (
  username,
  firstName,
  lastName,
  emailAddress,
  password,
  phoneNo,
  userType
) => {
  try {
    console.log("api-call-section");
    console.log(MAIN_URL + REGISTER_END_POINT);

    let response = await axios.post(MAIN_URL + REGISTER_END_POINT, {
      first_name: firstName,
      last_name: lastName,
      user_name: username,
      password: password,
      email: emailAddress,
      phone: phoneNo,
      user_type: userType
    });

    console.log(response);

    if(response.status===200 || response.status===201){
        console.log("success sign up ");
        return true;
    }

    console.log("response is not 200");

    return false;
   
  } catch (err) {
    console.log("failure in sign up");

    console.error(err);
    return false;
}
};
