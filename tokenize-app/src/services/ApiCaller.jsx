import axios from "axios";
import {
  MAIN_URL,
  LOGIN_END_POINT,
  REGISTER_END_POINT,
  LOG_OUT,
  DELETE_USER,
  DASH_BOARD,
  JWT,CSRF,CSRF_LABEL
} from "../Constants/config";


//create separate axios instance where u can configure default headers.
const axiosInstance = axios.create({
  baseURL:MAIN_URL,
  // withCredentials:true,
});

axiosInstance.interceptors.request.use(
  (config)=>{
    const csrfToken = sessionStorage.getItem(CSRF);
    console.log("Stored token : ",csrfToken);
    console.log("headers before ",config.headers);

    if (csrfToken){
      config.headers[CSRF_LABEL]=csrfToken;
    }

    console.log("headers after ",config.headers);
    return config;

  },(error)=>{
    if(error.response.status===401){
      console.log("Unauthorized 401 error !");
      console.error(error) 
    }
    return Promise.reject(error)
  }
);

export const loginApiCall = async (email_, password_) => {
  try {
    console.log("api-call-section");
    console.log(MAIN_URL + LOGIN_END_POINT);

    let response = await axiosInstance.post(LOGIN_END_POINT, {
      password: password_,
      email: email_,
    });

    console.log(response);

    if (response.status === 200 || response.status === 201) {
      console.log("success-login");
      sessionStorage.setItem(CSRF,response['data']['csrf']);
      return true;
    }

    console.log("login - failed!");
    return true;
  } catch (error) {
    console.log("faileure login");
    console.error(error);
    return false;
  }
};

export const navigateToDashboardApiCall = async () => {
  try {
    let response = await axiosInstance.get(DASH_BOARD)
    .then((result)=>{
      console.log("Success in navigation ",result);
    }).catch((err)=>{
        console.log("Error in navigation ",err)
    })

    console.log(response);

    if (response.status === 200 || response.status === 201) {
      console.log("success navigation !");
      return true;
    }
    console.log("navigation to dashboard status was not 2xx");
    return false;
  } catch (err) {
    console.log("something went wrong while navigating to dashboard");
    console.error(err);
    return false;
  }
};

export const deleteUserApiCall = async () => {
  console.log("delete user api invoked");
  console.log(MAIN_URL + DELETE_USER);

  try {
    let response = await axiosInstance.get( DELETE_USER);
    console.log(response)
    if (response.status === 200 || response.status === 201) {
      console.log("status is 2xx");
      sessionStorage.clear();
      return true;
    }

    console.log("status is not 2xx");
    return false;
  } catch (err) {
    console.error(err);
    return false;
  }
};

export const logoutApiCall = async () => {
  console.log("logout api call invoked");
  console.log(MAIN_URL + LOG_OUT);

  try {
    let response = await axiosInstance.get(LOG_OUT);
    console.log(response);

    if (response.status === 201 || response.status === 200) {
      console.log("status is 2xx");
      sessionStorage.clear();
      return true;
    }

    console.log("status was not 2xx");
    return false;
  } catch (err) {
    console.log("failed to logout");
    console.error(err);
    return false;
  }
};

export const signUpApiCall = async (
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

    let response = await axiosInstance.post(REGISTER_END_POINT, {
      first_name: firstName,
      last_name: lastName,
      user_name: username,
      password: password,
      email: emailAddress,
      phone: phoneNo,
      user_type: userType,
    });

    console.log(response);

    if (response.status === 200 || response.status === 201) {
      console.log("success sign up ");
      console.log(response['data']['csrf'])
      sessionStorage.setItem(CSRF,response['data']['csrf']);
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
