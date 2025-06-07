import * as requestService from "../services/requestServices";
import { checkCode, isArray } from "../utils/myUtils";
import { message } from "antd";

export default {
  namespace: "labelManagement",
  state: {
    labels: [],
    loading: false
  },
  reducers: {
    save(state, { payload }) {
      return { ...state, ...payload };
    },
    setLoading(state, { payload }) {
      return { ...state, loading: payload };
    }
  },
  effects: {
    *fetchLabels({ payload }, { call, put }) {
      try {
        yield put({ type: "setLoading", payload: true });
        const res = yield call(requestService.getAllQuestionLabels);
        if (checkCode(res) && isArray(res.data)) {
          yield put({
            type: "save",
            payload: {
              labels: res.data
            }
          });
        }
      } catch (e) {
        console.error("获取标签列表失败:", e);
        message.error("获取标签列表失败");
      } finally {
        yield put({ type: "setLoading", payload: false });
      }
    }
  }
};
