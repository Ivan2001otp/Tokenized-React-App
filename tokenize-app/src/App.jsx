
import { BrowserRouter as Router,Route,Routes } from 'react-router-dom';
import { Login } from './pages/Login'
import { SignUp } from './pages/SignUp'
import {Home} from './pages/Home';
import { ErrorPage } from './pages/ErrorPage';
import { Dashboard } from './pages/Dashboard';

function App() {

  return (
    <Router>
      <Routes>
      <Route path="/" element={<Home/>}/>
        <Route path="/login" element={<Login/>}/>
        <Route path="/signup" element={<SignUp/>}/>
        <Route path="/dashboard" element={<Dashboard/>}/>
        <Route path="/error-page" element={<ErrorPage/>}/>
      </Routes>
    </Router>
  );
}

export default App
