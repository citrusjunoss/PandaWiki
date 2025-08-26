import { updateAppDetail } from '@/api';
import ShowText from '@/components/ShowText';
import {
  Box,
  FormControlLabel,
  Radio,
  RadioGroup,
  TextField,
} from '@mui/material';
import { Message } from 'ct-mui';
import { useEffect, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import {
  DomainKnowledgeBaseDetail,
  DomainAppDetailResp,
} from '@/request/types';
import { getApiV1AppDetail } from '@/request/App';
import { FormItem, SettingCardItem } from './Common';

const CardRobotWecomService = ({
  kb,
  url,
}: {
  kb: DomainKnowledgeBaseDetail;
  url: string;
}) => {
  const [isEdit, setIsEdit] = useState(false);
  const [detail, setDetail] = useState<DomainAppDetailResp | null>(null);
  const [isEnabled, setIsEnabled] = useState(false);

  const {
    control,
    handleSubmit,
    formState: { errors },
    reset,
  } = useForm({
    defaultValues: {
      wechat_service_is_enabled: false,
      wechat_service_secret: '',
      wechat_service_token: '',
      wechat_service_encodingaeskey: '',
      wechat_service_corpid: '',
    },
  });

  const getDetail = () => {
    getApiV1AppDetail({ kb_id: kb.id!, type: '6' }).then(res => {
      setDetail(res);
      setIsEnabled(res.settings?.wechat_service_is_enabled ?? false);
      reset({
        wechat_service_is_enabled:
          res.settings?.wechat_service_is_enabled ?? false,
        wechat_service_secret: res.settings?.wechat_service_secret ?? '',
        wechat_service_token: res.settings?.wechat_service_token ?? '',
        wechat_service_encodingaeskey:
          res.settings?.wechat_service_encodingaeskey ?? '',
        wechat_service_corpid: res.settings?.wechat_service_corpid ?? '',
      });
    });
  };

  const onSubmit = handleSubmit(data => {
    if (!detail) return;
    updateAppDetail(
      { id: detail.id! },
      {
        settings: {
          wechat_service_is_enabled: data.wechat_service_is_enabled,
          wechat_service_secret: data.wechat_service_secret,
          wechat_service_token: data.wechat_service_token,
          wechat_service_encodingaeskey: data.wechat_service_encodingaeskey,
          wechat_service_corpid: data.wechat_service_corpid,
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
      title='企业微信客服'
      isEdit={isEdit}
      onSubmit={onSubmit}
      more={{
        type: 'link',
        href: 'https://pandawiki.docs.baizhi.cloud/node/01980888-bb0c-77d6-bc45-4636249b0d96',
        target: '_blank',
        text: '使用方法',
      }}
    >
      <FormItem label='企业微信客服'>
        <Controller
          control={control}
          name='wechat_service_is_enabled'
          render={({ field }) => (
            <RadioGroup
              row
              {...field}
              onChange={e => {
                field.onChange(e.target.value === 'true');
                setIsEnabled(e.target.value === 'true');
                setIsEdit(true);
              }}
            >
              <FormControlLabel
                value={true}
                control={<Radio size='small' />}
                label={<Box sx={{ width: 100 }}>启用</Box>}
              />
              <FormControlLabel
                value={false}
                control={<Radio size='small' />}
                label={<Box sx={{ width: 100 }}>禁用</Box>}
              />
            </RadioGroup>
          )}
        />
      </FormItem>

      {isEnabled && (
        <>
          <FormItem label='回调地址'>
            <ShowText text={[`${url}/share/v1/app/wechat/service`]} />
          </FormItem>
          <FormItem label='企业 ID' required>
            <Controller
              control={control}
              name='wechat_service_corpid'
              rules={{
                required: '企业 ID',
              }}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  placeholder=''
                  onChange={e => {
                    field.onChange(e.target.value);
                    setIsEdit(true);
                  }}
                  error={!!errors.wechat_service_corpid}
                  helperText={errors.wechat_service_corpid?.message}
                />
              )}
            />
          </FormItem>
          <FormItem label='Corp Secret' required>
            <Controller
              control={control}
              name='wechat_service_secret'
              rules={{
                required: 'Corp Secret',
              }}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  placeholder=''
                  onChange={e => {
                    field.onChange(e.target.value);
                    setIsEdit(true);
                  }}
                  error={!!errors.wechat_service_secret}
                  helperText={errors.wechat_service_secret?.message}
                />
              )}
            />
          </FormItem>
          <FormItem label='Token' required>
            <Controller
              control={control}
              name='wechat_service_token'
              rules={{
                required: 'Suite Token',
              }}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  placeholder=''
                  onChange={e => {
                    field.onChange(e.target.value);
                    setIsEdit(true);
                  }}
                  error={!!errors.wechat_service_token}
                  helperText={errors.wechat_service_token?.message}
                />
              )}
            />
          </FormItem>
          <FormItem label='Encoding Aes Key' required>
            <Controller
              control={control}
              name='wechat_service_encodingaeskey'
              rules={{
                required: 'Suite Encoding Aes Key',
              }}
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  placeholder=''
                  onChange={e => {
                    field.onChange(e.target.value);
                    setIsEdit(true);
                  }}
                  error={!!errors.wechat_service_encodingaeskey}
                  helperText={errors.wechat_service_encodingaeskey?.message}
                />
              )}
            />
          </FormItem>
        </>
      )}
    </SettingCardItem>
  );
};

export default CardRobotWecomService;
