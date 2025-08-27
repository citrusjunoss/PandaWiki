import { DomainKnowledgeBaseDetail, V1UserInfoResp } from '@/request/types';
import { KnowledgeBaseListItem } from '@/api';
import { DomainAppDetailResp } from '@/request/types';
import { createSlice } from '@reduxjs/toolkit';

export interface config {
  user: V1UserInfoResp;
  kb_id: string;
  kbList: KnowledgeBaseListItem[];
  kb_c: boolean;
  modelStatus: boolean;
  kbDetail: DomainKnowledgeBaseDetail;
  appPreviewData: DomainAppDetailResp | null;
}
const initialState: config = {
  user: {
    id: '',
    account: '',
    created_at: '',
  },
  kb_id: '',
  kbList: [],
  kb_c: false,
  modelStatus: false,
  kbDetail: {} as DomainKnowledgeBaseDetail,
  appPreviewData: null,
};

const configSlice = createSlice({
  name: 'config',
  initialState,
  reducers: {
    setUser(state, { payload }) {
      state.user = payload;
    },
    setKbId(state, { payload }) {
      localStorage.setItem('kb_id', payload);
      state.kb_id = payload;
    },
    setKbList(state, { payload }) {
      state.kbList = payload;
    },
    setKbC(state, { payload }) {
      state.kb_c = payload;
    },
    setModelStatus(state, { payload }) {
      state.modelStatus = payload;
    },
    setAppPreviewData(state, { payload }) {
      state.appPreviewData = payload;
    },
    setKbDetail(state, { payload }) {
      state.kbDetail = payload;
    },
  },
});

export const {
  setUser,
  setKbId,
  setKbList,
  setKbC,
  setModelStatus,
  setAppPreviewData,
  setKbDetail,
} = configSlice.actions;
export default configSlice.reducer;
