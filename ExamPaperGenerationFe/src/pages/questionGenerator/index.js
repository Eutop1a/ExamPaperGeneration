import React from "react";
import { connect, history } from "umi";
import {
  Button,
  Card,
  Descriptions,
  Divider,
  Input,
  InputNumber,
  PageHeader,
  Popover,
  Radio,
  Select,
  Spin,
  Switch,
  Tag,
  Tooltip,
  message,
  Modal,
  Form,
  Badge,
  Empty,
  Image
} from "antd";
import {
  QuestionCircleOutlined,
  StopOutlined,
  EyeOutlined,
  CheckCircleOutlined,
  PlusOutlined,
  EditOutlined,
  DeleteOutlined,
  SettingOutlined
} from "@ant-design/icons";
import style from "./index.less";
import { myEmptyStatus } from "../../layouts/commonComponents";
// echarts 按需加载
import "echarts/lib/chart/pie";
import "echarts/lib/component/tooltip";
import "echarts/lib/component/title";
import "echarts/lib/component/legend";
import "echarts/lib/component/markPoint";
import ReactEcharts from "echarts-for-react";
import {
  getAllQuestionLabels,
  createLabel,
  updateLabel,
  deleteLabel,
  setSimilarityThreshold
} from "../../services/requestServices";

class questionGenerator extends React.Component {
  constructor(props) {
    super(props);
    this.state = {
      // 按钮loading
      btnLoading: false,
      genBtnLoading: false,
      // 组卷参数设置
      testPaperName: "试卷_" + new Date().getTime(),
      // 随机选择开关
      randomSwitch: true,
      // 遗传算法
      geneticSelect: 0,
      // 迭代次数
      iterationsNum: 400,
      TKTCount: 10,
      XZTCount: 10,
      PDTCount: 5,
      JDTCount: 2,
      // 平均难度
      averageDifficulty: 2.75,
      // 知识点范围
      chapter1Range: [],
      generateRange: [],
      // 知识点权重
      knowledgeWeights: [],
      // 题目类型要求
      questionTypeRequirements: {
        选择题: {
          minCount: 5,
          targetScore: 25
        },
        填空题: {
          minCount: 5,
          targetScore: 25
        },
        判断题: {
          minCount: 5,
          targetScore: 25
        },
        简答题: {
          minCount: 2,
          targetScore: 25
        }
      },
      similarityModalVisible: false,
      similarityThreshold: 0.8
    };
    this.similarityFormRef = React.createRef();
  }

  // btn handler
  // 自动组卷！
  autoGenerate = async () => {
    // 检查知识点权重总和是否为100%
    const totalWeight = this.calculateTotalWeight();
    if (totalWeight !== 100) {
      message.error(`知识点权重总和必须为100%，当前为${totalWeight}%`);
      return;
    }

    await this.setState({ btnLoading: true });
    await this.props.dispatch({ type: "questionGenerator/emptyList" });
    let selectedTopicIds = [];
    await this.props.testPaperGenList.forEach(item => {
      selectedTopicIds.push(item.id);
    });

    let payload = {
      selectedTopicIds: selectedTopicIds,
      generateRange: this.state.generateRange,
      averageDifficulty: this.state.averageDifficulty,
      knowledgeWeights: this.state.knowledgeWeights,
      questionTypeRequirements: this.state.questionTypeRequirements
    };

    if (this.state.geneticSelect === 0) {
      await this.props.dispatch({
        type: "questionGenerator/randomSelect",
        payload: payload
      });
    } else if (this.state.geneticSelect === 1) {
      payload["iterationsNum"] = this.state.iterationsNum;
      try {
        const response = await this.props.dispatch({
          type: "questionGenerator/geneticSelect",
          payload: payload
        });
        console.log("Genetic algorithm response:", response);
        if (response && response.data) {
          // 确保数据被正确更新到state中
          await this.props.dispatch({
            type: "questionGenerator/updateGeneticResults",
            payload: response.data
          });
        }
      } catch (error) {
        console.error("Error in genetic algorithm:", error);
        message.error("遗传算法执行出错，请重试");
      }
    }
    await this.setState({ btnLoading: false });
  };
  // 生成试卷按钮
  handleClickGenTestPaper = async () => {
    await this.setState({ genBtnLoading: true });
    let questionIdList = [];
    await this.props.testPaperGenList.forEach(item => {
      questionIdList.push(item.id);
    });
    await this.props.dispatch({
      type: "questionGenerator/questionGen2",
      payload: {
        questionIdList: questionIdList,
        randomSwitch: this.state.randomSwitch,
        testPaperName: this.state.testPaperName
      }
    });
    await this.setState({ genBtnLoading: false });
  };

  // data calc
  // 计算所有题目的总数
  getAllTopicCount = () => {
    const testPaperGenList = this.props.testPaperGenList || [];
    const TKTList = this.props.TKTList || [];
    const PDTList = this.props.PDTList || [];
    const JDTList = this.props.JDTList || [];
    const XZTList = this.props.XZTList || [];

    return (
      testPaperGenList.length +
      TKTList.length +
      PDTList.length +
      JDTList.length +
      XZTList.length
    );
  };
  // 计算自动出题的总数
  getRandomSelectTopicCount = () => {
    const TKTList = this.props.TKTList || [];
    const PDTList = this.props.PDTList || [];
    const JDTList = this.props.JDTList || [];
    const XZTList = this.props.XZTList || [];

    return TKTList.length + PDTList.length + JDTList.length + XZTList.length;
  };
  // 计算总平均难度
  calcAllAvgDifficulty = () => {
    const testPaperGenList = this.props.testPaperGenList || [];
    const TKTList = this.props.TKTList || [];
    const PDTList = this.props.PDTList || [];
    const JDTList = this.props.JDTList || [];
    const XZTList = this.props.XZTList || [];

    const _totalCount =
      testPaperGenList.length +
      TKTList.length +
      PDTList.length +
      JDTList.length +
      XZTList.length;
    if (_totalCount === 0) return "无";
    let _sumDifficulty = 0.0;
    testPaperGenList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    TKTList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    PDTList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    JDTList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    XZTList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    return (_sumDifficulty / _totalCount).toFixed(2);
  };
  // 计算手动选择的题目难度
  calcManualSelectTopicDifficulty = () => {
    const testPaperGenList = this.props.testPaperGenList || [];
    if (testPaperGenList.length === 0) return "无";
    let _sumDifficulty = 0.0;
    testPaperGenList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    return (_sumDifficulty / testPaperGenList.length).toFixed(2);
  };
  // 计算自动出题平均难度
  calcRandomSelectTopicDifficulty = () => {
    const TKTList = this.props.TKTList || [];
    const PDTList = this.props.PDTList || [];
    const JDTList = this.props.JDTList || [];
    const XZTList = this.props.XZTList || [];

    const _totalCount =
      TKTList.length + PDTList.length + JDTList.length + XZTList.length;
    if (_totalCount === 0) return "无";
    let _sumDifficulty = 0.0;
    TKTList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    PDTList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    JDTList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    XZTList.forEach(item => {
      _sumDifficulty += item.difficulty;
    });
    return (_sumDifficulty / _totalCount).toFixed(2);
  };

  // 生成echarts的配置项
  getEchartsOption = () => {
    const testPaperGenList = this.props.testPaperGenList || [];
    const TKTList = this.props.TKTList || [];
    const PDTList = this.props.PDTList || [];
    const JDTList = this.props.JDTList || [];
    const XZTList = this.props.XZTList || [];

    let TKTCount = TKTList.length;
    let XZTCount = PDTList.length;
    let PDTCount = JDTList.length;
    let JDTCount = XZTList.length;

    testPaperGenList.forEach(item => {
      if (item.topic_type === "填空题") TKTCount++;
      else if (item.topic_type === "选择题") XZTCount++;
      else if (item.topic_type === "判断题") PDTCount++;
      else if (item.topic_type === "简答题") JDTCount++;
    });

    return {
      tooltip: {
        trigger: "item",
        formatter: "{b}数量 : {c} ({d}%)"
      },
      series: [
        {
          name: "数量",
          type: "pie",
          radius: [0, "50%"],
          label: {
            position: "outside",
            fontSize: 10
          },
          labelLine: {
            length: 10,
            length2: 15,
            smooth: true
          },
          data: [
            { value: TKTCount, name: "填空题" },
            { value: XZTCount, name: "判断题" },
            { value: PDTCount, name: "选择题" },
            { value: JDTCount, name: "简答题" }
          ],
          emphasis: {
            itemStyle: {
              shadowBlur: 10,
              shadowOffsetX: 0,
              shadowColor: "rgba(0, 0, 0, 0.5)"
            }
          }
        }
      ]
    };
  };

  getVarianceEchartsOption = () => {
    const variance = this.props.variance || [];

    return {
      xAxis: {
        type: "category",
        name: "迭代次数"
      },
      yAxis: {
        type: "value",
        name: "与预期难度的方差"
      },
      tooltip: {
        trigger: "axis",
        formatter: "遗传迭代次数：{b}<br />当前方差： {c}"
      },
      series: [
        {
          data: variance,
          type: "line",
          lineStyle: {
            color: {
              type: "linear",
              x: 0,
              y: 0,
              x2: 0,
              y2: 1,
              colorStops: [
                {
                  offset: 0,
                  color: "red" // 0% 处的颜色
                },
                {
                  offset: 1,
                  color: "blue" // 100% 处的颜色
                }
              ],
              global: false // 缺省为 false
            }
          },
          smooth: true
        }
      ]
    };
  };

  // 计算知识点权重总和
  calculateTotalWeight = () => {
    return this.state.knowledgeWeights.reduce(
      (sum, item) => sum + item.weight,
      0
    );
  };

  // init Data
  initData = async () => {
    await this.props.dispatch({ type: "questionEdit/getAllQuestionLabels" });
    this.setState({
      chapter1Range: this.props.label1,
      generateRange: this.props.label1
    });
    // 初始化知识点权重
    this.initKnowledgeWeights();
  };

  // 初始化知识点权重
  initKnowledgeWeights = () => {
    const { label1 } = this.props;
    if (!label1 || label1.length === 0) return;

    // 计算每个知识点的默认权重，确保总和为100
    const defaultWeight = Math.floor(100 / label1.length);
    const remainder = 100 - defaultWeight * label1.length;

    const weights = label1.map((label, index) => ({
      label1: label,
      weight: index < remainder ? defaultWeight + 1 : defaultWeight
    }));

    this.setState({ knowledgeWeights: weights });
  };

  // 更新知识点权重
  updateKnowledgeWeight = (index, value) => {
    if (value === null) return;

    const newWeights = [...this.state.knowledgeWeights];
    newWeights[index].weight = value;

    this.setState({ knowledgeWeights: newWeights });
  };

  // lifeCycle
  componentDidMount() {
    this.initData().then(null);
  }

  // 处理相似度阈值设置
  handleSimilaritySetting = () => {
    this.setState({ similarityModalVisible: true });
  };

  handleSimilaritySubmit = async () => {
    try {
      const values = await this.similarityFormRef.current.validateFields();
      console.log("设置相似度阈值:", values.threshold);
      await setSimilarityThreshold(values.threshold);
      message.success("相似度阈值设置成功");
      this.setState({ similarityModalVisible: false });
    } catch (error) {
      console.error("设置相似度阈值失败:", error);
      message.error(`设置失败: ${error.message || "未知错误"}`);
    }
  };

  render() {
    const renderManualTopic = () => {
      const testPaperGenList = this.props.testPaperGenList || [];
      if (testPaperGenList.length > 0) {
        return testPaperGenList.map((item, index) => {
          return (
            <Descriptions span={1} key={index}>
              <Descriptions.Item>
                {index + 1}、(本题{item.score}分) {item.topic}
              </Descriptions.Item>
              {item.topic_image_path &&
                item.topic_image_path.startsWith("data:image") && (
                  <Descriptions.Item>
                    <Image
                      src={item.topic_image_path}
                      alt="题目图片"
                      style={{ maxWidth: "300px", maxHeight: "200px" }}
                      preview={true}
                    />
                  </Descriptions.Item>
                )}
            </Descriptions>
          );
        });
      } else return myEmptyStatus("无数据", "200px");
    };

    const renderRandomTopic = () => {
      const testPaperGenList = this.props.testPaperGenList || [];
      const TKTList = this.props.TKTList || [];
      const PDTList = this.props.PDTList || [];
      const JDTList = this.props.JDTList || [];
      const XZTList = this.props.XZTList || [];
      const variance = this.props.variance || [];

      console.log("Variance in render:", variance);
      console.log("Chart options:", this.getVarianceEchartsOption());

      return (
        <div>
          {this.state.geneticSelect === 1 && variance.length > 0 ? (
            <div>
              <Divider
                orientation="left"
                style={{ fontWeight: "bold", fontSize: "1em" }}
              >
                遗传迭代算法 「自动出题难度」与「预期难度」的方差变化
              </Divider>
              <div style={{ marginBottom: "10px" }}>
                <code
                  style={{
                    color: "#666",
                    fontSize: "0.8em",
                    paddingRight: "15px"
                  }}
                >
                  初始方差：{variance[0].toFixed(4)}
                </code>
                <code
                  style={{
                    color: "#666",
                    fontSize: "0.8em",
                    paddingRight: "15px"
                  }}
                >
                  最终方差：{variance[variance.length - 1].toFixed(4)}
                </code>
                <code
                  style={{
                    color: "#666",
                    fontSize: "0.8em",
                    paddingRight: "15px"
                  }}
                >
                  迭代次数：{variance.length}
                </code>
              </div>
              <ReactEcharts
                option={this.getVarianceEchartsOption()}
                id="variance"
                style={{
                  width: "100%",
                  height: "250px",
                  transition: "all 1s",
                  border: "1px solid #f0f0f0",
                  borderRadius: "4px",
                  padding: "10px"
                }}
              />
            </div>
          ) : null}
          <Divider orientation="left" style={{ fontWeight: "bold" }}>
            选择题，共{XZTList.length}题
          </Divider>
          {XZTList.length > 0
            ? XZTList.map((item, index) => {
                return (
                  <Descriptions span={1} key={index}>
                    <Descriptions.Item>
                      {index + 1}、(本题{item.score}分) {item.topic}
                    </Descriptions.Item>
                    {item.topic_image_path &&
                      item.topic_image_path.startsWith("data:image") && (
                        <Descriptions.Item>
                          <Image
                            src={item.topic_image_path}
                            alt="题目图片"
                            style={{ maxWidth: "300px", maxHeight: "200px" }}
                            preview={true}
                          />
                        </Descriptions.Item>
                      )}
                  </Descriptions>
                );
              })
            : myEmptyStatus("无数据", "200px")}
          <Divider orientation="left" style={{ fontWeight: "bold" }}>
            填空题，共{TKTList.length}题
          </Divider>
          {TKTList.length > 0
            ? TKTList.map((item, index) => {
                return (
                  <Descriptions span={1} key={index}>
                    <Descriptions.Item>
                      {index + 1}、(本题{item.score}分) {item.topic}
                    </Descriptions.Item>
                    {item.topic_image_path &&
                      item.topic_image_path.startsWith("data:image") && (
                        <Descriptions.Item>
                          <Image
                            src={item.topic_image_path}
                            alt="题目图片"
                            style={{ maxWidth: "300px", maxHeight: "200px" }}
                            preview={true}
                          />
                        </Descriptions.Item>
                      )}
                  </Descriptions>
                );
              })
            : myEmptyStatus("无数据", "200px")}
          <Divider orientation="left" style={{ fontWeight: "bold" }}>
            判断题，共{PDTList.length}题
          </Divider>
          {PDTList.length > 0
            ? PDTList.map((item, index) => {
                return (
                  <Descriptions span={1} key={index}>
                    <Descriptions.Item>
                      {index + 1}、(本题{item.score}分) {item.topic}
                    </Descriptions.Item>
                    {item.topic_image_path &&
                      item.topic_image_path.startsWith("data:image") && (
                        <Descriptions.Item>
                          <Image
                            src={item.topic_image_path}
                            alt="题目图片"
                            style={{ maxWidth: "300px", maxHeight: "200px" }}
                            preview={true}
                          />
                        </Descriptions.Item>
                      )}
                  </Descriptions>
                );
              })
            : myEmptyStatus("无数据", "200px")}
          <Divider orientation="left" style={{ fontWeight: "bold" }}>
            简答题，共{JDTList.length}题
          </Divider>
          {JDTList.length > 0
            ? JDTList.map((item, index) => {
                return (
                  <Descriptions span={1} key={index}>
                    <Descriptions.Item>
                      {index + 1}、(本题{item.score}分) {item.topic}
                    </Descriptions.Item>
                    {item.topic_image_path &&
                      item.topic_image_path.startsWith("data:image") && (
                        <Descriptions.Item>
                          <Image
                            src={item.topic_image_path}
                            alt="题目图片"
                            style={{ maxWidth: "300px", maxHeight: "200px" }}
                            preview={true}
                          />
                        </Descriptions.Item>
                      )}
                  </Descriptions>
                );
              })
            : myEmptyStatus("无数据", "200px")}
        </div>
      );
    };

    const renderSummary = () => {
      const testPaperGenList = this.props.testPaperGenList || [];
      const TKTList = this.props.TKTList || [];
      const PDTList = this.props.PDTList || [];
      const JDTList = this.props.JDTList || [];
      const XZTList = this.props.XZTList || [];

      return (
        <div>
          <div className={style.middle_line_space}>
            <span>题目总数：</span>
            <Tag>{this.getAllTopicCount()}</Tag>
          </div>
          <div className={style.middle_line_space}>
            <span>平均难度：</span>
            <Tag>{this.calcAllAvgDifficulty()}</Tag>
            <Popover content="显示手动选择的题目和自动出题的平均难度">
              <QuestionCircleOutlined />
            </Popover>
          </div>
          <div className={style.middle_line_space}>
            <span>手动选择：</span>
            <Tag>{testPaperGenList.length}</Tag>
          </div>
          <div className={style.middle_line_space}>
            <span>手动选择的题目难度：</span>
            <Tag>{this.calcManualSelectTopicDifficulty()}</Tag>
            <Popover content="显示手动选择题目的平均难度">
              <QuestionCircleOutlined />
            </Popover>
          </div>
          <div className={style.middle_line_space}>
            <span>自动出题数：</span>
            <Tag>{this.getRandomSelectTopicCount()}</Tag>
          </div>
          <div className={style.middle_line_space}>
            <span>自动出题难度：</span>
            <Tag>{this.calcRandomSelectTopicDifficulty()}</Tag>
            <Popover
              title="为什么与预设值不相符？"
              content="题库中没有相应的符合预设要求的题目，会自动选取与预设值相近的题目"
            >
              <QuestionCircleOutlined />
            </Popover>
          </div>
          <Divider />
          <ReactEcharts
            option={this.getEchartsOption()}
            id="summary_echarts"
            style={{
              width: "100%",
              height: this.getAllTopicCount() > 0 ? "180px" : "0",
              transition: "all 1s"
            }}
          />
          <span>输入试卷名称（备注）：</span>
          <Input
            onChange={e => this.setState({ testPaperName: e.target.value })}
            value={this.state.testPaperName}
          />
          <Divider />
          <Button
            icon={<CheckCircleOutlined />}
            onClick={this.handleClickGenTestPaper}
            type="primary"
            style={{ display: "block", margin: "5px auto" }}
            loading={this.state.genBtnLoading}
            disabled={
              testPaperGenList.length +
                TKTList.length +
                PDTList.length +
                JDTList.length +
                XZTList.length <=
              0
            }
          >
            立即生成试卷word文档
          </Button>
        </div>
      );
    };

    const renderActionBox = () => {
      return (
        <div>
          <Spin
            spinning={!this.state.randomSwitch}
            indicator={<StopOutlined />}
            tip="关闭自动组卷功能，当前手动组卷"
          >
            <div className={style.text_line_space}>
              <span>试卷随机出题总数：</span>
              <Tag>
                {this.state.generateRange.length > 0
                  ? "根据知识点权重自动分配"
                  : "请先选择知识点"}
              </Tag>
            </div>
            <Divider />
            <div className={style.middle_line_space}>设置总体难度：</div>
            <InputNumber
              type="number"
              min={1}
              max={5}
              value={this.state.averageDifficulty}
              onChange={value =>
                value ? this.setState({ averageDifficulty: value }) : null
              }
              className={style.wrapper_params_input}
            />
            <div className={style.middle_line_space}>设置出题知识点：</div>
            <Select
              mode="multiple"
              style={{ width: "100%" }}
              placeholder="请选择出题范围"
              value={this.state.generateRange}
              onChange={value => {
                this.setState({ generateRange: value });
                // 更新知识点权重
                if (value && value.length > 0) {
                  const weights = value.map(label => {
                    return {
                      label1: label,
                      weight: Math.floor(100 / value.length) // 默认平均分配权重，向下取整
                    };
                  });
                  // 处理余数，将余数加到最后一个权重上
                  const remainder =
                    100 - weights.reduce((sum, item) => sum + item.weight, 0);
                  if (remainder > 0 && weights.length > 0) {
                    weights[weights.length - 1].weight += remainder;
                  }
                  this.setState({ knowledgeWeights: weights });
                } else {
                  this.setState({ knowledgeWeights: [] });
                }
              }}
            >
              {this.state.chapter1Range.map((item, index) => (
                <Select.Option key={index} value={item}>
                  {item}
                </Select.Option>
              ))}
            </Select>

            {this.state.generateRange.length > 0 && (
              <>
                <Divider />
                <div className={style.middle_line_space}>设置知识点权重：</div>
                {this.state.knowledgeWeights.map((item, index) => (
                  <div key={index} className={style.middle_line_space}>
                    <span>{item.label1}：</span>
                    <InputNumber
                      type="number"
                      min={0}
                      max={100}
                      value={item.weight}
                      onChange={value =>
                        this.updateKnowledgeWeight(index, value)
                      }
                      className={style.wrapper_params_input}
                      style={{ width: "70px" }}
                    />
                    <span style={{ marginLeft: "5px" }}>%</span>
                  </div>
                ))}
                <div
                  className={style.middle_line_space}
                  style={{ marginTop: "10px" }}
                >
                  <span>权重总和：</span>
                  <Tag
                    color={
                      this.calculateTotalWeight() === 100 ? "green" : "red"
                    }
                  >
                    {this.calculateTotalWeight()}%
                  </Tag>
                  {this.calculateTotalWeight() !== 100 && (
                    <span style={{ marginLeft: "10px", color: "red" }}>
                      (权重总和必须等于100%)
                    </span>
                  )}
                </div>
                <div
                  className={style.middle_line_space}
                  style={{ marginTop: "10px" }}
                >
                  <Button
                    size="small"
                    onClick={() => {
                      const weights = this.state.knowledgeWeights.map(item => {
                        return {
                          ...item,
                          weight: Math.floor(
                            100 / this.state.knowledgeWeights.length
                          )
                        };
                      });
                      // 处理余数，将余数加到最后一个权重上
                      const remainder =
                        100 -
                        weights.reduce((sum, item) => sum + item.weight, 0);
                      if (remainder > 0 && weights.length > 0) {
                        weights[weights.length - 1].weight += remainder;
                      }
                      this.setState({ knowledgeWeights: weights });
                    }}
                  >
                    平均分配权重
                  </Button>
                </div>
              </>
            )}

            <Divider />
            <div className={style.middle_line_space}>设置题目类型要求：</div>
            {Object.entries(this.state.questionTypeRequirements).map(
              ([type, requirements], index) => (
                <div
                  key={index}
                  className={style.middle_line_space}
                  style={{ display: "flex", alignItems: "center", gap: "10px" }}
                >
                  <span style={{ minWidth: "60px" }}>{type}</span>
                  <div
                    style={{
                      display: "flex",
                      alignItems: "center",
                      gap: "10px",
                      flex: 1
                    }}
                  >
                    <span>最少数量：</span>
                    <InputNumber
                      type="number"
                      min={0}
                      max={100}
                      value={requirements.minCount}
                      onChange={value => {
                        if (value !== null) {
                          const newRequirements = {
                            ...this.state.questionTypeRequirements
                          };
                          newRequirements[type].minCount = value;
                          this.setState({
                            questionTypeRequirements: newRequirements
                          });
                        }
                      }}
                      style={{ width: "70px" }}
                    />
                    <span>目标分数：</span>
                    <InputNumber
                      type="number"
                      min={0}
                      max={100}
                      value={requirements.targetScore}
                      onChange={value => {
                        if (value !== null) {
                          const newRequirements = {
                            ...this.state.questionTypeRequirements
                          };
                          newRequirements[type].targetScore = value;
                          this.setState({
                            questionTypeRequirements: newRequirements
                          });
                        }
                      }}
                      style={{ width: "70px" }}
                    />
                  </div>
                </div>
              )
            )}

            <Divider />
            <Radio.Group
              onChange={e => this.setState({ geneticSelect: e.target.value })}
              value={this.state.geneticSelect}
            >
              <Radio value={0}>贪心最优算法</Radio>
              <Radio value={1}>遗传迭代算法</Radio>
            </Radio.Group>
            <div className={style.middle_line_space}>设置遗传迭代次数：</div>
            <InputNumber
              type="number"
              min={0}
              max={10000}
              value={this.state.iterationsNum}
              disabled={this.state.geneticSelect !== 1}
              onChange={value =>
                value ? this.setState({ iterationsNum: value }) : null
              }
              className={style.wrapper_params_input}
            />
            <Divider />
            <Tooltip
              placement="topLeft"
              title="结果会显示在左侧【自动组卷】面板"
            >
              <Button
                icon={<CheckCircleOutlined />}
                onClick={this.autoGenerate}
                type="primary"
                loading={this.state.btnLoading}
                style={{ display: "block", margin: "5px auto", width: "130px" }}
              >
                自动组卷！
              </Button>
            </Tooltip>
          </Spin>
        </div>
      );
    };

    return (
      <div>
        <PageHeader
          title={"试题组卷" + "（当前共" + this.getAllTopicCount() + "题）"}
          subTitle={"预览即将生成的试题，或调整随机抽题参数"}
          extra={[
            <Button
              type="primary"
              icon={<EyeOutlined />}
              onClick={() => {
                history.push("/questionGenHistory");
              }}
              key="1"
            >
              查看出题历史
            </Button>,
            <Tooltip title="设置题目相似度阈值" key="similarity">
              <Button
                type="primary"
                icon={<SettingOutlined />}
                onClick={this.handleSimilaritySetting}
              >
                相似度设置
              </Button>
            </Tooltip>
          ]}
        />
        <div className={style.flex_wrapper}>
          {/*左侧题目预览*/}
          {/*手动出题预览*/}
          <div className={style.wrapper_left_side}>
            <Card
              title={
                <span className={style.preview_title}>
                  手动出题，当前已手动选择{this.props.testPaperGenList.length}题
                </span>
              }
            >
              {renderManualTopic()}
            </Card>
            <Divider />
            {/*自动组卷预览*/}
            {this.state.randomSwitch ? (
              <Card
                title={
                  <span className={style.preview_title}>
                    自动组卷{" "}
                    {this.state.geneticSelect === 1 ? (
                      <Tag color="magenta">遗传迭代算法</Tag>
                    ) : (
                      <Tag color="geekblue">贪心最优算法</Tag>
                    )}
                  </span>
                }
              >
                {renderRandomTopic()}
              </Card>
            ) : null}
          </div>
          {/*右侧自动组卷功能*/}
          <div className={style.wrapper_right_side}>
            <Card
              title={
                <div className={style.preview_title}>
                  <span style={{ marginRight: "20px" }}>自动组卷功能</span>
                  <Switch
                    checked={this.state.randomSwitch}
                    onChange={checked =>
                      this.setState({ randomSwitch: checked })
                    }
                    checkedChildren="开自动组卷"
                    unCheckedChildren="关自动组卷"
                  />
                </div>
              }
            >
              {renderActionBox()}
            </Card>
          </div>
          {/*右侧生成试卷功能*/}
          <div className={style.wrapper_right_side}>
            <Card
              title={
                <span className={style.preview_title}>这份试卷的统计信息</span>
              }
            >
              {renderSummary()}
            </Card>
          </div>
        </div>

        {/* 相似度阈值设置模态框 */}
        <Modal
          title="设置题目相似度阈值"
          visible={this.state.similarityModalVisible}
          onOk={this.handleSimilaritySubmit}
          onCancel={() => this.setState({ similarityModalVisible: false })}
        >
          <Form ref={this.similarityFormRef} layout="vertical">
            <Form.Item
              name="threshold"
              label="相似度阈值"
              initialValue={0.8}
              rules={[
                { required: true, message: "请输入相似度阈值" },
                {
                  type: "number",
                  min: 0,
                  max: 1,
                  message: "阈值必须在0到1之间"
                }
              ]}
            >
              <InputNumber
                style={{ width: "100%" }}
                step={0.1}
                precision={2}
                placeholder="请输入0-1之间的数值"
              />
            </Form.Item>
          </Form>
        </Modal>
      </div>
    );
  }
}

function mapStateToProps({ questionGenerator, questionEdit }) {
  const {
    testPaperGenList,
    TKTList,
    PDTList,
    JDTList,
    XZTList,
    variance
  } = questionGenerator;
  const {
    allQuestionLabels,
    chapter1,
    chapter2,
    label1,
    label2
  } = questionEdit;
  return {
    testPaperGenList,
    TKTList,
    PDTList,
    JDTList,
    XZTList,
    variance,
    allQuestionLabels,
    chapter1,
    chapter2,
    label1,
    label2
  };
}

export default connect(mapStateToProps)(questionGenerator);
