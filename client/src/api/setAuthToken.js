import Api from './api';

const setAuthToken = (token) => {
  if (token) {
    Api.defaults.headers.common['Authorization'] = 'Bearer ' + token;
  } else {
    delete Api.defaults.headers.common['Authorization'];
  }
};

export default setAuthToken;
