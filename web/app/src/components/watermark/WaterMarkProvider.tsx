'use client';
import React from 'react';
import dayjs from 'dayjs';
import Watermark, { WatermarkProps } from './index';
import { useStore } from '@/provider';

const WaterMarkProvider = (props: WatermarkProps) => {
  const { children, ...rest } = props;
  const { kbDetail, authInfo } = useStore();
  const content = kbDetail?.settings?.watermark_content;
  const enable = kbDetail?.settings?.watermark_enable;
  if (!enable) {
    return children;
  }
  const time = `${authInfo?.username ?? ''} ${dayjs().format('YYYY-MM-DD HH:mm:ss')}`;
  const contentLines = [time, ...(content?.split('\n') || [])];
  return (
    <Watermark {...rest} content={contentLines}>
      {children}
    </Watermark>
  );
};

export default WaterMarkProvider;
