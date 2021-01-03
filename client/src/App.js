import './App.css';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import Root from './components/screens/Root';
import Login from './components/screens/Login';
import Register from './components/screens/Register';
import Home from './components/screens/Home';
import Navbar from './components/layout/Navbar';
import setAuthToken from './api/setAuthToken';
import AuthState from './context/auth/AuthState';

if (localStorage.token) {
  setAuthToken(localStorage.token);
}

const App = () => {
  return (
    <AuthState>
      <Router>
        <Navbar />
        <Switch>
          <Route exact path="/" component={Root} />
          <Route exact path="/register" component={Register} />
          <Route exact path="/login" component={Login} />
          <Route exact path="/home" component={Home} />
        </Switch>
      </Router>
    </AuthState>
  );
};

export default App;
