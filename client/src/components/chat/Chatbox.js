import { Menu } from 'antd';
import React, { useState } from 'react';
import { MessageOutlined, UserOutlined } from '@ant-design/icons';
import FriendsList from './FriendsList';
import MessagesList from './MessagesList';
import ChatTaskbar from './ChatTaskbar';

const Chatbox = () => {
  const [current, setCurrent] = useState(null);

  const handleClick = (e) => {
    setCurrent(e.key);
  };

  return (
    <>
      <Menu onClick={handleClick} selectedKeys={[current]} mode="horizontal">
        <Menu.Item key="chat" icon={<MessageOutlined />}>
          Messages
        </Menu.Item>
        <Menu.Item key="friends" icon={<UserOutlined />}>
          Friends
        </Menu.Item>
      </Menu>
      {current === 'chat' ? <MessagesList /> : <FriendsList />}
      <ChatTaskbar />
    </>
  );
};

export default Chatbox;
