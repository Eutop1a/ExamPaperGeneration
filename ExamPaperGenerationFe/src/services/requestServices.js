import { request } from "umi";
import { API } from "../config/requestConfig";
import { parseParams } from "../utils/myUtils";
import axios from "axios";
import { message } from "antd";

// 创建axios实例
const requestService = axios.create({
  baseURL: API,
  timeout: 10000,
  withCredentials: true
});

// 响应拦截器
requestService.interceptors.response.use(
  response => {
    return response.data;
  },
  error => {
    console.error("请求错误:", error);
    message.error(`请求失败: ${error.message || "未知错误"}`);
    return Promise.reject(error);
  }
);

// 检查响应状态码
const checkCode = response => {
  if (response.code === 200) {
    return response.data;
  } else {
    message.error(response.message || "操作失败");
    return Promise.reject(new Error(response.message || "操作失败"));
  }
};

// 导出requestService和checkCode
export { requestService, checkCode };

// login
export function login(payload) {
  const url = `${API}/login`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

// getLoginStatus
export function getLoginStatus() {
  const url = `${API}/getLoginStatus`;
  return request(url, {
    method: "get",
    mode: "cors",
    credentials: "include"
  });
}

// logout
export function logout() {
  const url = `${API}/logout`;
  return request(url, {
    method: "post",
    mode: "cors",
    credentials: "include"
  });
}

// registered
export function registered(payload) {
  const url = `${API}/registered`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

// getAllUser
export function getAllUser() {
  const url = `${API}/getAllUser`;
  return request(url, {
    method: "get",
    mode: "cors",
    credentials: "include"
  });
}

// deleteUser
export function deleteUser(payload) {
  const url = `${API}/deleteUser`;
  return request(url, {
    method: "get",
    mode: "cors",
    params: payload,
    credentials: "include"
  });
}

// getApplyUser
export function getApplyUser() {
  const url = `${API}/getApplyUser`;
  return request(url, {
    method: "get",
    mode: "cors",
    credentials: "include"
  });
}

// passApply
export function passApply(payload) {
  const url = `${API}/passApply`;
  return request(url, {
    method: "get",
    mode: "cors",
    params: payload,
    credentials: "include"
  });
}

// deleteApply
export function deleteApply(payload) {
  const url = `${API}/deleteApply`;
  return request(url, {
    method: "get",
    mode: "cors",
    params: payload,
    credentials: "include"
  });
}

// questionBank

export function getQuestionBank() {
  const url = `${API}/getAllQuestionBank`;
  return request(url, {
    method: "get",
    mode: "cors",
    credentials: "include"
  });
}

export function deleteSingleQuestionBank(payload) {
  const url = `${API}/deleteSingleQuestionBank`;
  return request(url, {
    method: "get",
    params: payload,
    mode: "cors",
    credentials: "include"
  });
}

export function getEachChapterCount() {
  const url = `${API}/getEachChapterCount`;
  return request(url, {
    method: "get",
    mode: "cors",
    credentials: "include"
  });
}

export function getEachScoreCount() {
  const url = `${API}/getEachScoreCount`;
  return request(url, {
    method: "get",
    mode: "cors",
    credentials: "include"
  });
}

// questionEdit

export function getAllQuestionLabels() {
  const url = `${API}/getAllQuestionLabels`;
  return request(url, {
    method: "get",
    mode: "cors",
    credentials: "include"
  });
}

export function insertSingleQuestionBank(payload) {
  const url = `${API}/insertSingleQuestionBank`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

export function getQuestionBankById(payload) {
  const url = `${API}/getQuestionBankById`;
  return request(url, {
    method: "get",
    params: payload,
    mode: "cors",
    credentials: "include"
  });
}

export function updateQuestionBankById(payload) {
  const url = `${API}/updateQuestionBankById`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

// questionGenerator
export function questionGen(payload) {
  const url = `${API}/questionGen`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

export function randomSelect(payload) {
  const url = `${API}/randomSelect`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

export function geneticSelect(payload) {
  const url = `${API}/geneticSelect`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

export function downloadFile() {
  // 文件下载
  const a = document.createElement("a");
  document.body.append(a);
  const url = `${API}/getFile`;
  a.href = url;
  a.download = "试卷.docx";
  a.target = "_blank";
  a.click();
  a.remove();
  window.URL.revokeObjectURL(url);
}

// questionGenHistory
export function getAllTestPaperGenHistory() {
  const url = `${API}/getAllTestPaperGenHistory`;
  return request(url, {
    method: "get",
    mode: "cors",
    credentials: "include"
  });
}

export function getQuestionGenHistoriesByTestPaperUid(payload) {
  const url = `${API}/getQuestionGenHistoriesByTestPaperUid`;
  return request(url, {
    method: "get",
    params: payload,
    mode: "cors",
    credentials: "include"
  });
}

export function deleteQuestionGenHistoryByTestPaperUid(payload) {
  const url = `${API}/deleteQuestionGenHistoryByTestPaperUid`;
  return request(url, {
    method: "get",
    params: payload,
    mode: "cors",
    credentials: "include"
  });
}

export function uploadFile(payload) {
  const url = `${API}/upload`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

// 新 生成word接口
export function questionGen2(payload) {
  const url = `${API}/questionGen2`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include",
    parseResponse: false // 关闭返回值简化
  });
}

// 重新导出word接口
export function reExportTestPaper(payload) {
  const url = `${API}/reExportTestPaper`;
  return request(url, {
    method: "get",
    params: payload,
    mode: "cors",
    credentials: "include",
    parseResponse: false // 关闭返回值简化
  });
}

// 导出答案接口
export function exportAnswer(payload) {
  const url = `${API}/exportAnswer`;
  return request(url, {
    method: "get",
    params: payload,
    mode: "cors",
    credentials: "include",
    parseResponse: false // 关闭返回值简化
  });
}

// 获取所有题目类型
export function getTopicType(payload) {
  const url = `${API}/getTopicType`;
  return request(url, {
    method: "get",
    params: payload,
    mode: "cors",
    credentials: "include"
  });
}

// 搜索题目
export function searchQuestionByTopic(payload) {
  const url = `${API}/searchQuestionByTopic`;
  return request(url, {
    method: "get",
    params: payload,
    mode: "cors",
    credentials: "include"
  });
}

// 新 更新已经组卷的题目
export function updateQuestionGenHistory(payload) {
  const url = `${API}/updateQuestionGenHistory`;
  return request(url, {
    method: "post",
    data: payload,
    mode: "cors",
    credentials: "include"
  });
}

// 标签管理相关接口
export const getLabels = () => {
  return request(`${API}/labels`, {
    method: "GET",
    mode: "cors",
    credentials: "include"
  });
};

export const createLabel = data => {
  return request(`${API}/labels`, {
    method: "POST",
    mode: "cors",
    credentials: "include",
    data
  });
};

export const updateLabel = (id, data) => {
  return request(`${API}/labels/${id}`, {
    method: "PUT",
    mode: "cors",
    credentials: "include",
    data
  });
};

export const deleteLabel = id => {
  return request(`${API}/labels/${id}`, {
    method: "DELETE",
    mode: "cors",
    credentials: "include"
  });
};

// 相似度阈值设置
export const setSimilarityThreshold = threshold => {
  return request(`${API}/similarity`, {
    method: "PUT",
    mode: "cors",
    credentials: "include",
    data: { threshold }
  });
};

// 更新题目（支持图片上传）
export const updateQuestionWithImage = async (id, formData) => {
  try {
    const response = await requestService({
      url: "/updateQuestionBankById",
      method: "post",
      data: formData,
      headers: {
        "Content-Type": "multipart/form-data"
      }
    });
    return checkCode(response);
  } catch (error) {
    console.error("更新题目失败:", error);
    throw error;
  }
};

// 创建题目（支持图片上传）
export const insertQuestionWithImage = async formData => {
  try {
    const response = await requestService({
      url: "/insertSingleQuestionBankWithImg",
      method: "post",
      data: formData,
      headers: {
        "Content-Type": "multipart/form-data"
      }
    });
    return checkCode(response);
  } catch (error) {
    console.error("创建题目失败:", error);
    throw error;
  }
};
