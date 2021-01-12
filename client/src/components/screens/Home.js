import React, { useState, useContext, useEffect } from 'react';
import AuthContext from '../../context/auth/authContext';
import 'antd/dist/antd.css';
import { Button, Drawer } from 'antd';
import Chatbox from '../chat/Chatbox';

const Home = () => {
  const authContext = useContext(AuthContext);
  const [chat, setChat] = useState(false);

  useEffect(
    () => {
      authContext.loadUser();
    },
    // eslint-disable-next-line
    []
  );

  const showChat = () => {
    setChat(true);
  };

  const hideChat = () => {
    setChat(false);
  };

  return (
    <>
      <Button
        type="primary"
        size="large"
        style={{ position: 'absolute', right: '5%', bottom: '5%' }}
        onClick={showChat}
      >
        WAHOOOO!
      </Button>
      <Drawer visible={chat} onClose={hideChat} width={350}>
        <Chatbox />
      </Drawer>
    </>
  );
};

export default Home;
