import React from "react";
import { connect, history } from "umi";
import {
  Alert,
  Button,
  Form,
  Input,
  Radio,
  Modal,
  Select,
  InputNumber,
  DatePicker,
  message,
  Statistic,
  Divider
} from "antd";
import { LockOutlined, LoginOutlined, UserOutlined } from "@ant-design/icons";
import styles from "./index.less";
import moment from "moment";
import ReactEcharts from "echarts-for-react";
import { delay } from "../../utils/myUtils";

class OverViewModal extends React.Component {
  constructor(props) {
    super(props);
    this.state = {};
  }

  onCancel = async () => {
    this.props.hide();
  };

  onOk = () => {
    this.onCancel().then();
  };

  // calc
  getEachTopicTypeCount = () => {
    const totalCount = this.props.dataSource.length;
    let TKTCount = 0;
    let XZTCount = 0;
    let PDTCount = 0;
    let JDTCount = 0;
    this.props.dataSource.forEach((item, index) => {
      if (item.topic_type === "填空题") TKTCount++;
      else if (item.topic_type === "选择题") XZTCount++;
      else if (item.topic_type === "判断题") PDTCount++;
      else if (item.topic_type === "简答题") JDTCount++;
    });
    return { totalCount, TKTCount, XZTCount, PDTCount, JDTCount };
  };

  echarts_option_1_pie_chart = () => {
    let data = [];
    if (this.props.eachChapterCount) {
      this.props.eachChapterCount.forEach((item, index) => {
        data.push({ name: item.label_1, value: item.count });
      });
    }
    return {
      title: {
        text: "各章节知识点分布（饼图）"
      },
      tooltip: {
        trigger: "item",
        formatter: "{b} : {c}题"
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
          data: data,
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

  echarts_option_2_pie_chart = () => {
    const result = this.getEachTopicTypeCount();
    const data = [];
    data.push({ name: "填空题", value: result.TKTCount });
    data.push({ name: "选择题", value: result.XZTCount });
    data.push({ name: "判断题", value: result.PDTCount });
    data.push({ name: "简答题", value: result.JDTCount });
    return {
      title: {
        text: "题目类型分布",
        left: "center"
      },
      tooltip: {
        trigger: "item",
        formatter: "{a} <br/>{b}: {c} ({d}%)"
      },
      legend: {
        orient: "vertical",
        left: 10,
        data: ["填空题", "选择题", "判断题", "简答题"]
      },
      series: [
        {
          name: "题目类型",
          type: "pie",
          radius: ["50%", "70%"],
          avoidLabelOverlap: false,
          label: {
            show: false,
            position: "center"
          },
          emphasis: {
            label: {
              show: true,
              fontSize: "30",
              fontWeight: "bold"
            }
          },
          labelLine: {
            show: false
          },
          data: data
        }
      ]
    };
  };

  echarts_option_3_bar_chart = () => {
    let xAxisData = [];
    let seriesData = [];
    if (this.props.eachScoreCount) {
      this.props.eachScoreCount.forEach((item, index) => {
        xAxisData.push(item.score);
        seriesData.push(item.count);
      });
    }
    return {
      title: {
        text: "各分值分布统计（柱状图）"
      },
      xAxis: {
        type: "category",
        data: xAxisData,
        name: "每题分值"
      },
      yAxis: {
        type: "value",
        name: "题目数"
      },
      tooltip: {
        trigger: "axis",
        formatter: "{b}分题的数量 : {c}题"
      },
      series: [
        {
          data: seriesData,
          type: "bar",
          showBackground: false,
          itemStyle: {
            color: "#1890ff"
          },
          backgroundStyle: {
            color: "rgba(220, 220, 220, 0.8)"
          }
        }
      ]
    };
  };

  // initData
  initData = async () => {
    await this.props.dispatch({ type: "questionBank/getEachChapterCount" });
    await this.props.dispatch({ type: "questionBank/getEachScoreCount" });
  };

  // lifeCycle
  componentDidMount() {
    this.initData().then(null);
  }

  render() {
    const renderStat = () => {
      const totalCount = this.props.dataSource.length;
      let TKTCount = 0;
      let XZTCount = 0;
      let PDTCount = 0;
      let JDTCount = 0;
      this.props.dataSource.forEach((item, index) => {
        if (item.topic_type === "填空题") TKTCount++;
        else if (item.topic_type === "选择题") XZTCount++;
        else if (item.topic_type === "判断题") PDTCount++;
        else if (item.topic_type === "简答题") JDTCount++;
      });
      return (
        <div>
          <div className={styles.row_flex}>
            <Statistic
              className={styles.statistic_board}
              title="试题库总数："
              value={Number(totalCount)}
              suffix={"题"}
            />
            <Divider type="vertical" />
            <Statistic
              className={styles.statistic_board}
              title="填空题总数："
              value={Number(TKTCount)}
              suffix={"题"}
            />
            <Divider type="vertical" />
            <Statistic
              className={styles.statistic_board}
              title="选择题总数："
              value={Number(XZTCount)}
              suffix={"题"}
            />
            <Divider type="vertical" />
            <Statistic
              className={styles.statistic_board}
              title="判断题总数："
              value={Number(PDTCount)}
              suffix={"题"}
            />
            <Divider type="vertical" />
            <Statistic
              className={styles.statistic_board}
              title="简答题总数："
              value={Number(JDTCount)}
              suffix={"题"}
            />
            <Divider type="vertical" />
          </div>
          <Divider orientation="left" style={{ fontWeight: "bold" }}>
            题库各类数量统计
          </Divider>
          <div className={styles.row_flex}>
            <div style={{ width: "600px", height: "300px" }}>
              <ReactEcharts option={this.echarts_option_1_pie_chart()} />
            </div>
          </div>
          <div className={styles.row_flex}>
            <div style={{ width: "600px", height: "300px" }}>
              <ReactEcharts option={this.echarts_option_2_pie_chart()} />
            </div>
          </div>
          <div className={styles.row_flex}>
            <div style={{ width: "600px", height: "300px" }}>
              <ReactEcharts option={this.echarts_option_3_bar_chart()} />
            </div>
          </div>
        </div>
      );
    };

    return (
      <div>
        <Modal
          title="试题库概览"
          centered={true}
          visible={this.props.visible}
          onOk={this.onOk}
          onCancel={this.onCancel}
          footer={null}
          width={"70%"}
        >
          <div>{renderStat()}</div>
        </Modal>
      </div>
    );
  }
}

function mapStateToProps({ questionBank }) {
  const { eachChapterCount, eachScoreCount } = questionBank;
  return { eachChapterCount, eachScoreCount };
}

export default connect(mapStateToProps)(OverViewModal);
