import Api from './api';

const setAuthToken = (token) => {
  if (token) {
    Api.defaults.headers.common['authorization'] = 'Bearer ' + token;
  } else {
    delete Api.defaults.headers.common['athorization'];
  }
};

export default setAuthToken;
