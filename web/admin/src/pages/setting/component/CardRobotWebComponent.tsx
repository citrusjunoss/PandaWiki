import { updateAppDetail } from '@/api';
import { DomainKnowledgeBaseDetail } from '@/request/types';
import ShowText from '@/components/ShowText';
import UploadFile from '@/components/UploadFile';
import {
  Box,
  FormControlLabel,
  Link,
  Radio,
  RadioGroup,
  Stack,
  TextField,
} from '@mui/material';
import { Icon, Message } from 'ct-mui';
import { useEffect, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { FormItem, SettingCardItem } from './Common';
import { DomainAppDetailResp } from '@/request/types';
import { getApiV1AppDetail } from '@/request/App';

interface CardRobotWebComponentProps {
  kb: DomainKnowledgeBaseDetail;
}

const CardRobotWebComponent = ({ kb }: CardRobotWebComponentProps) => {
  const [isEdit, setIsEdit] = useState(false);
  const [isEnabled, setIsEnabled] = useState(false);
  const [detail, setDetail] = useState<DomainAppDetailResp | null>(null);
  const {
    control,
    handleSubmit,
    reset,
    formState: { errors },
  } = useForm({
    defaultValues: {
      is_open: 0,
      theme_mode: 'light',
      btn_text: '',
      btn_logo: '',
    },
  });

  const [url, setUrl] = useState<string>('');

  useEffect(() => {
    if (kb.access_settings?.base_url) {
      setUrl(kb.access_settings.base_url);
      return;
    }
    const host = kb.access_settings?.hosts?.[0] || '';
    if (host === '') return;
    const { ssl_ports = [], ports = [] } = kb.access_settings || {};

    if (ssl_ports) {
      if (ssl_ports.includes(443)) setUrl(`https://${host}`);
      else if (ssl_ports.length > 0) setUrl(`https://${host}:${ssl_ports[0]}`);
    } else if (ports) {
      if (ports.includes(80)) setUrl(`http://${host}`);
      else if (ports.length > 0) setUrl(`http://${host}:${ports[0]}`);
    }
  }, [kb]);

  const getDetail = () => {
    getApiV1AppDetail({ kb_id: kb.id!, type: '2' }).then(res => {
      setDetail(res);
      reset({
        is_open: res.settings?.widget_bot_settings?.is_open ? 1 : 0,
        theme_mode: res.settings?.widget_bot_settings?.theme_mode || 'light',
        btn_text: res.settings?.widget_bot_settings?.btn_text || '在线客服',
        btn_logo: res.settings?.widget_bot_settings?.btn_logo,
      });
      setIsEnabled(res.settings?.widget_bot_settings?.is_open ? true : false);
    });
  };

  const onSubmit = handleSubmit(data => {
    if (!detail) return;
    updateAppDetail(
      { id: detail.id! },
      {
        settings: {
          widget_bot_settings: {
            is_open: data.is_open === 1 ? true : false,
            theme_mode: data.theme_mode as 'light' | 'dark',
            btn_text: data.btn_text,
            btn_logo: data.btn_logo,
          },
        },
      },
    ).then(() => {
      Message.success('保存成功');
      setIsEdit(false);
      getDetail();
      reset();
    });
  });

  useEffect(() => {
    getDetail();
  }, [kb]);

  return (
    <SettingCardItem
      title='网页挂件机器人'
      isEdit={isEdit}
      onSubmit={onSubmit}
      more={
        <Link
          component='a'
          href='https://pandawiki.docs.baizhi.cloud/node/0197f335-a1a8-786c-95df-0848f61fb98a'
          target='_blank'
          sx={{
            fontSize: 14,
            textDecoration: 'none',
            fontWeight: 'normal',
            ml: 1,
            '&:hover': {
              fontWeight: 'bold',
            },
          }}
        >
          使用方法
        </Link>
      }
    >
      <FormItem label='网页挂件机器人'>
        <Controller
          control={control}
          name='is_open'
          render={({ field }) => (
            <RadioGroup
              row
              {...field}
              onChange={e => {
                field.onChange(+e.target.value as 1 | 0);
                setIsEnabled((+e.target.value as 1 | 0) === 1);
                setIsEdit(true);
              }}
            >
              <FormControlLabel
                value={1}
                control={<Radio size='small' />}
                label={<Box sx={{ width: 100 }}>启用</Box>}
              />
              <FormControlLabel
                value={0}
                control={<Radio size='small' />}
                label={<Box sx={{ width: 100 }}>禁用</Box>}
              />
            </RadioGroup>
          )}
        />
      </FormItem>
      {isEnabled && (
        <>
          <FormItem label='配色方案'>
            <Controller
              control={control}
              name='theme_mode'
              render={({ field }) => (
                <RadioGroup
                  row
                  {...field}
                  onChange={e => {
                    field.onChange(e.target.value);
                    setIsEdit(true);
                  }}
                >
                  <FormControlLabel
                    value={'light'}
                    control={<Radio size='small' />}
                    label={<Box sx={{ width: 100 }}>浅色模式</Box>}
                  />
                  <FormControlLabel
                    value={'dark'}
                    control={<Radio size='small' />}
                    label={<Box sx={{ width: 100 }}>深色模式</Box>}
                  />
                </RadioGroup>
              )}
            />
          </FormItem>
          <FormItem label='侧边按钮文字'>
            <Controller
              control={control}
              name='btn_text'
              render={({ field }) => (
                <TextField
                  fullWidth
                  {...field}
                  placeholder='输入侧边按钮文字'
                  error={!!errors.btn_text}
                  helperText={errors.btn_text?.message}
                  onChange={event => {
                    setIsEdit(true);
                    field.onChange(event);
                  }}
                />
              )}
            />
          </FormItem>
          <FormItem label='侧边按钮 Logo'>
            <Controller
              control={control}
              name='btn_logo'
              render={({ field }) => (
                <UploadFile
                  {...field}
                  id='btn_logo'
                  type='url'
                  accept='image/*'
                  width={80}
                  onChange={url => {
                    field.onChange(url);
                    setIsEdit(true);
                  }}
                />
              )}
            />
          </FormItem>
          <FormItem label='嵌入代码'>
            {url ? (
              <ShowText
                noEllipsis
                text={[
                  `<!--// Head 标签引入样式 -->`,
                  `<link rel="stylesheet" href="${url}/widget-bot.css">`,
                  '',
                  `<!--// Body 标签引入挂件 -->`,
                  `<script src="${url}/widget-bot.js"></script>`,
                ]}
              />
            ) : (
              <Stack
                direction='row'
                alignItems={'center'}
                gap={0.5}
                sx={{ color: 'warning.main', fontSize: 14 }}
              >
                <Icon type='icon-jinggao' />
                未配置域名，可在右侧
                <Box component={'span'} sx={{ fontWeight: 500 }}>
                  服务监听方式
                </Box>{' '}
                中配置
              </Stack>
            )}
          </FormItem>
        </>
      )}
    </SettingCardItem>
  );
};

export default CardRobotWebComponent;
