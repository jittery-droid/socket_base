import './App.css';
import { BrowserRouter as Router, Switch, Route } from 'react-router-dom';
import Root from './components/screens/Root';
import Login from './components/screens/Login';
import Register from './components/screens/Register';
import Navbar from './components/layout/Navbar';

const App = () => {
  return (
    <Router>
      <Navbar />
      <Switch>
        <Route exact path="/" component={Root} />
        <Route exact path="/register" component={Register} />
        <Route exact path="/login" component={Login} />
      </Switch>
    </Router>
  );
};

export default App;
