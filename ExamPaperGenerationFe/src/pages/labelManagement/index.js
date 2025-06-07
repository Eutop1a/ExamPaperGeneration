import React from "react";
import { connect, history } from "umi";
import {
  Table,
  Button,
  Modal,
  Form,
  Input,
  message,
  Popconfirm,
  Space,
  Card,
  PageHeader
} from "antd";
import { PlusOutlined, EditOutlined, DeleteOutlined } from "@ant-design/icons";
import { renderLoading } from "../../layouts/commonComponents";
import style from "./index.less";
import {
  getAllQuestionLabels,
  createLabel,
  updateLabel,
  deleteLabel
} from "../../services/requestServices";

class LabelManagement extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      isLoading: false,
      visible: false,
      editingLabel: null,
      error: null
    };
    this._isMounted = false;
    this.formRef = React.createRef();
  }

  // initData
  initData = async () => {
    if (!this._isMounted) return;

    try {
      console.log("开始加载标签数据...");
      this.setState({ isLoading: true, error: null });
      const response = await getAllQuestionLabels();
      console.log("标签数据响应:", response);

      if (this._isMounted) {
        if (response && response.data) {
          this.props.dispatch({
            type: "labelManagement/save",
            payload: {
              labels: response.data || []
            }
          });
          console.log("标签数据已更新到状态:", response.data);
        } else {
          console.warn("响应数据格式不正确:", response);
          throw new Error("响应数据格式不正确");
        }
        this.setState({ isLoading: false });
      }
    } catch (error) {
      console.error("加载标签数据失败:", error);
      if (this._isMounted) {
        this.setState({
          isLoading: false,
          error: error.message || "加载数据失败"
        });
        message.error(`加载数据失败: ${error.message || "未知错误"}`);
      }
    }
  };

  // lifeCycle
  componentDidMount() {
    console.log("组件已挂载");
    this._isMounted = true;
    this.initData();
  }

  componentWillUnmount() {
    console.log("组件即将卸载");
    this._isMounted = false;
  }

  // handlers
  handleCreate = () => {
    if (!this._isMounted) return;

    console.log("打开创建标签模态框");
    this.setState({
      editingLabel: null,
      visible: true
    });
    if (this.formRef.current) {
      this.formRef.current.resetFields();
    }
  };

  handleEdit = record => {
    if (!this._isMounted) return;

    console.log("打开编辑标签模态框:", record);
    this.setState({
      editingLabel: record,
      visible: true
    });
    if (this.formRef.current) {
      this.formRef.current.setFieldsValue(record);
    }
  };

  handleDelete = async id => {
    if (!this._isMounted) return;

    try {
      console.log("开始删除标签:", id);
      await deleteLabel(id);
      if (this._isMounted) {
        console.log("标签删除成功");
        message.success("删除成功");
        this.initData();
      }
    } catch (error) {
      console.error("删除标签失败:", error);
      if (this._isMounted) {
        message.error(`删除失败: ${error.message || "未知错误"}`);
      }
    }
  };

  handleSubmit = async () => {
    if (!this._isMounted) return;

    try {
      console.log("开始提交表单");
      const values = await this.formRef.current.validateFields();
      console.log("表单验证通过，提交数据:", values);

      if (this.state.editingLabel) {
        const { id } = this.state.editingLabel;
        console.log("更新标签:", { id, values });
        await updateLabel(id, values);
        if (this._isMounted) {
          console.log("标签更新成功");
          message.success("更新成功");
          this.setState({ visible: false });
          this.initData();
        }
      } else {
        console.log("创建新标签:", values);
        await createLabel(values);
        if (this._isMounted) {
          console.log("标签创建成功");
          message.success("创建成功");
          this.setState({ visible: false });
          this.initData();
        }
      }
    } catch (error) {
      console.error("提交表单失败:", error);
      if (this._isMounted) {
        message.error(`操作失败: ${error.message || "未知错误"}`);
      }
    }
  };

  render() {
    console.log("渲染组件，当前状态:", {
      isLoading: this.state.isLoading,
      error: this.state.error,
      labels: this.props.labels,
      loading: this.props.loading
    });

    if (this.state.error) {
      return (
        <div className={style.wrapper}>
          <PageHeader
            title="知识点标签管理"
            subTitle="管理试题的知识点标签"
            style={{ background: "#fff", marginBottom: 16 }}
          />
          <Card>
            <div style={{ color: "red", padding: "20px" }}>
              加载失败: {this.state.error}
              <Button
                type="primary"
                onClick={this.initData}
                style={{ marginLeft: "10px" }}
              >
                重试
              </Button>
            </div>
          </Card>
        </div>
      );
    }

    if (this.state.isLoading) return renderLoading("正在加载", "50vh");

    const columns = [
      {
        title: "章节1",
        dataIndex: "chapter_1",
        key: "chapter_1"
      },
      {
        title: "章节2",
        dataIndex: "chapter_2",
        key: "chapter_2"
      },
      {
        title: "标签1",
        dataIndex: "label_1",
        key: "label_1"
      },
      {
        title: "标签2",
        dataIndex: "label_2",
        key: "label_2"
      },
      {
        title: "操作",
        key: "action",
        render: (_, record) => (
          <Space size="middle">
            <Button
              type="text"
              icon={<EditOutlined />}
              onClick={() => this.handleEdit(record)}
            >
              编辑
            </Button>
            <Popconfirm
              title="确定要删除这个标签吗？"
              onConfirm={() => this.handleDelete(record.id)}
              okText="确定"
              cancelText="取消"
            >
              <Button type="text" danger icon={<DeleteOutlined />}>
                删除
              </Button>
            </Popconfirm>
          </Space>
        )
      }
    ];

    return (
      <div className={style.wrapper}>
        <PageHeader
          title="知识点标签管理"
          subTitle="管理试题的知识点标签"
          style={{ background: "#fff", marginBottom: 16 }}
        />
        <Card>
          <div style={{ marginBottom: 16 }}>
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={this.handleCreate}
            >
              新建标签
            </Button>
          </div>
          <Table
            columns={columns}
            dataSource={this.props.labels}
            rowKey="id"
            loading={this.props.loading}
          />
          <Modal
            title={this.state.editingLabel ? "编辑标签" : "新建标签"}
            visible={this.state.visible}
            onOk={this.handleSubmit}
            onCancel={() => this.setState({ visible: false })}
          >
            <Form ref={this.formRef} layout="vertical">
              <Form.Item
                name="chapter_1"
                label="章节1"
                rules={[{ required: true, message: "请输入章节1" }]}
              >
                <Input />
              </Form.Item>
              <Form.Item
                name="chapter_2"
                label="章节2"
                rules={[{ required: true, message: "请输入章节2" }]}
              >
                <Input />
              </Form.Item>
              <Form.Item
                name="label_1"
                label="标签1"
                rules={[{ required: true, message: "请输入标签1" }]}
              >
                <Input />
              </Form.Item>
              <Form.Item
                name="label_2"
                label="标签2"
                rules={[{ required: true, message: "请输入标签2" }]}
              >
                <Input />
              </Form.Item>
            </Form>
          </Modal>
        </Card>
      </div>
    );
  }
}

function mapStateToProps({ labelManagement }) {
  console.log("mapStateToProps - labelManagement state:", labelManagement);
  return {
    labels: labelManagement.labels,
    loading: labelManagement.loading
  };
}

export default connect(mapStateToProps)(LabelManagement);
