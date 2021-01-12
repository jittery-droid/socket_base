import { Button, Input } from 'antd';
import React, { useState } from 'react';
import { UserAddOutlined, UserDeleteOutlined } from '@ant-design/icons';
import { Modal, Form } from 'antd';
import Api from '../../api/api';

const ChatTaskbar = () => {
  const [addModal, showAddModal] = useState(false);
  const [deleteModal, showDeleteModal] = useState(false);
  const [form] = Form.useForm();

  const addFriend = async (friend) => {
    // move to friend context
    await Api.post(
      '/api/friends',
      { friend },
      {
        headers: {
          'Content-Type': 'application/json',
          'Access-Control-Allow-Origin': '*',
        },
      }
    );
  };

  const removeFriend = (friend) => {};

  return (
    <div
      style={{
        width: '20%',
        position: 'absolute',
        bottom: '5%',
        right: '5%',
        display: 'flex',
        justifyContent: 'space-between',
      }}
    >
      <Button
        type="primary"
        icon={<UserAddOutlined />}
        onClick={() => showAddModal(true)}
      ></Button>
      <Button
        type="primary"
        icon={<UserDeleteOutlined />}
        onClick={() => showDeleteModal(true)}
      ></Button>
      <Modal
        visible={addModal}
        onCancel={() => showAddModal(false)}
        title="Add Friend"
        maskClosable={true}
        onOk={() => {
          form.validateFields().then((values) => {
            form.resetFields();
            addFriend(values);
          });
        }}
      >
        <Form form={form}>
          <Form.Item name="name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
      <Modal
        visible={deleteModal}
        onCancel={() => showDeleteModal(false)}
        title="Remove Friend"
        maskClosable={true}
        onOk={() => {
          form.validateFields().then((values) => {
            form.resetFields();
            removeFriend(values);
          });
        }}
      >
        <Form form={form}>
          <Form.Item name="name" rules={[{ required: true }]}>
            <Input />
          </Form.Item>
        </Form>
      </Modal>
    </div>
  );
};

export default ChatTaskbar;
