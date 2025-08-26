import { useAppSelector } from '@/store';
import { Box, Stack, Typography } from '@mui/material';
import { ThemeProvider } from 'ct-mui';

import { Dispatch, SetStateAction, useMemo, useState, useEffect } from 'react';
import { AppSetting } from '@/api';
import Header from './Header';
import componentStyleOverrides from '@/themes/override';
import light from '../theme/light';
import dark from '../theme/dark';
import { Component } from '..';

interface ShowContentProps {
  curComponent: Component;
  setCurComponent: Dispatch<SetStateAction<string>>;
  renderMode: 'pc' | 'mobile';
  scale: number;
}

const ShowContent = ({
  setCurComponent,
  curComponent,
  renderMode,
  scale,
}: ShowContentProps) => {
  const { appPreviewData } = useAppSelector(state => state.config);
  const [key, setKey] = useState(0); // 用于强制重新渲染ThemeProvider的key

  // 监听主题模式变化，强制重新渲染ThemeProvider
  useEffect(() => {
    setKey(prev => prev + 1);
  }, [appPreviewData?.settings?.theme_mode]);

  // @ts-expect-error 类型错误
  const settings: Partial<AppSetting> = useMemo(() => {
    return (
      appPreviewData?.settings || {
        title: '默认标题',
        icon: '',
        btns: [],
        header_search_placeholder: '',
        allow_theme_switching: false,
      }
    );
  }, [appPreviewData]);

  // 渲染带高亮边框的组件
  const renderHighlightedComponent = (
    componentName: string,
    component: React.ReactNode,
  ) => {
    const isHighlighted = curComponent.name === componentName;

    return (
      <Box
        sx={{
          position: 'relative',
          border: isHighlighted ? '2px solid #5F58FE' : '2px solid transparent',
          borderRadius: '0px',
          padding: '2px',
          cursor: 'pointer',
          '&:hover': {
            border: isHighlighted ? '2px solid #5F58FE' : '2px dashed #5F58FE',
          },
        }}
        // 添加自定义属性用于标识组件
        data-component={componentName}
        onClick={e => {
          setCurComponent(componentName);
        }}
      >
        {component}
        {isHighlighted && (
          <Typography
            sx={{
              position: 'absolute',
              left: '-2px',
              bottom: '-24px',
              fontWeight: 400,
              lineHeight: '22px',
              bgcolor: '#5F58FE',
              color: '#FFFFFF',
              fontSize: '14px',
              padding: '1px 16px',
            }}
          >
            {curComponent.title}
          </Typography>
        )}
      </Box>
    );
  };

  return (
    <ThemeProvider
      key={key} // 强制刷新主题
      colors={{ light, dark }}
      mode={appPreviewData?.settings!.theme_mode || 'light'}
      theme={{
        components: componentStyleOverrides,
      }}
    >
      <Stack
        sx={{
          flex: 1,
          flexShrink: 0,
          height: '95%',
          marginTop: '20px',
          width: '100%',
          overflowX: 'auto',
          overflowY: 'auto',
          borderRight: '1px solid #ECEEF1',
          borderLeft: '1px solid #ECEEF1',
          borderTop: '1px solid #ECEEF1',
          '&::-webkit-scrollbar': {
            height: '8px', // 滚动条高度
          },
          '&::-webkit-scrollbar-thumb': {
            background: '#888', // 滑块颜色
            borderRadius: '4px',
          },
        }}
      >
        <Stack
          sx={{
            minWidth: renderMode === 'pc' ? `1200px` : '375px',
            width: renderMode === 'pc' ? `100%` : '375px',
            margin: '0 auto',
            boxShadow:
              renderMode === 'pc' ? null : '0 10px 15px -3px rgb(0 0 0 / 0.1)',
            height: '100%',
            overflowX: renderMode === 'pc' ? 'auto' : 'hidden',
            overflowY: 'auto',
            position: 'relative',
            bgcolor: 'background.default',
            transform: `scale(${scale})`,
            transformOrigin: 'left top',
            transition: 'transform 0.2s ease',
          }}
        >
          {/* Header预览部分 */}
          {renderHighlightedComponent(
            'header',
            <Header settings={settings} renderMode={renderMode} />,
          )}
        </Stack>
      </Stack>
    </ThemeProvider>
  );
};

export default ShowContent;
