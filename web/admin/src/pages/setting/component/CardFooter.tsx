import { updateAppDetail } from '@/api';
import { FooterSetting } from '@/api/type';
import DragBrand from '@/components/Drag/DragBrand';
import UploadFile from '@/components/UploadFile';
import { DomainAppDetailResp, DomainBrandGroup } from '@/request/types';
import {
  Box,
  Button,
  FormControlLabel,
  Radio,
  RadioGroup,
  Stack,
  TextField,
} from '@mui/material';
import { Message } from 'ct-mui';
import { useEffect, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { FormItem, SettingCardItem } from './Common';

interface CardFooterProps {
  id: string;
  data: DomainAppDetailResp;
  refresh: (value: FooterSetting) => void;
}

const CardFooter = ({ id, data, refresh }: CardFooterProps) => {
  const [isEdit, setIsEdit] = useState(false);
  const {
    control,
    handleSubmit,
    watch,
    setValue,
    formState: { errors },
  } = useForm({
    defaultValues: {
      footer_style: 'simple',
      corp_name: '',
      icp: '',
      brand_name: '',
      brand_desc: '',
      brand_logo: '',
      brand_groups: [] as DomainBrandGroup[],
    },
  });

  const footerStyle = watch('footer_style');

  const onSubmit = handleSubmit(value => {
    updateAppDetail(
      { id },
      {
        settings: {
          ...data.settings,
          // @ts-expect-error 类型不匹配
          footer_settings: {
            ...data.settings?.footer_settings,
            ...value,
          },
        },
      },
    ).then(() => {
      refresh(value as FooterSetting);
      Message.success('保存成功');
      setIsEdit(false);
    });
  });

  useEffect(() => {
    setValue(
      'footer_style',
      data.settings?.footer_settings?.footer_style as 'simple' | 'complex',
    );
    setValue('corp_name', data.settings?.footer_settings?.corp_name || '');
    setValue('icp', data.settings?.footer_settings?.icp || '');
    setValue('brand_name', data.settings?.footer_settings?.brand_name || '');
    setValue('brand_desc', data.settings?.footer_settings?.brand_desc || '');
    setValue('brand_logo', data.settings?.footer_settings?.brand_logo || '');
    setValue(
      'brand_groups',
      data.settings?.footer_settings?.brand_groups || [],
    );
  }, [data]);

  return (
    <SettingCardItem title='页脚' isEdit={isEdit} onSubmit={onSubmit}>
      <FormItem label='页面模式'>
        <Controller
          control={control}
          name='footer_style'
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
                value={'simple'}
                control={<Radio size='small' />}
                label={<Box sx={{ width: 100 }}>简单页脚</Box>}
              />
              <FormControlLabel
                value={'complex'}
                control={<Radio size='small' />}
                label={<Box sx={{ width: 100 }}>扩展页脚</Box>}
              />
            </RadioGroup>
          )}
        />
      </FormItem>
      <FormItem label='企业名称 / 组织名称'>
        {' '}
        <Controller
          control={control}
          name='corp_name'
          render={({ field }) => (
            <TextField
              {...field}
              fullWidth
              placeholder='企业名称/组织名称'
              onChange={e => {
                field.onChange(e.target.value);
                setIsEdit(true);
              }}
              error={!!errors.corp_name}
              helperText={errors.corp_name?.message}
            />
          )}
        />
      </FormItem>

      <FormItem label='ICP 备案编号'>
        <Controller
          control={control}
          name='icp'
          render={({ field }) => (
            <TextField
              {...field}
              fullWidth
              placeholder='ICP 备案编号'
              onChange={e => {
                field.onChange(e.target.value);
                setIsEdit(true);
              }}
              error={!!errors.icp}
              helperText={errors.icp?.message}
            />
          )}
        />
      </FormItem>

      {footerStyle === 'complex' && (
        <>
          <FormItem label='品牌名称'>
            <Controller
              control={control}
              name='brand_name'
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  placeholder='品牌名称'
                  onChange={e => {
                    field.onChange(e.target.value);
                    setIsEdit(true);
                  }}
                  error={!!errors.brand_name}
                  helperText={errors.brand_name?.message}
                />
              )}
            />
          </FormItem>

          <FormItem label='品牌 Logo'>
            <Controller
              control={control}
              name='brand_logo'
              render={({ field }) => (
                <UploadFile
                  {...field}
                  id='brand_logo'
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

          <FormItem label='品牌介绍'>
            <Controller
              control={control}
              name='brand_desc'
              render={({ field }) => (
                <TextField
                  {...field}
                  fullWidth
                  placeholder='品牌介绍'
                  onChange={e => {
                    field.onChange(e.target.value);
                    setIsEdit(true);
                  }}
                  error={!!errors.brand_desc}
                  helperText={errors.brand_desc?.message}
                />
              )}
            />
          </FormItem>

          <Box>
            <Box
              sx={{
                width: 156,
                fontSize: 14,
                lineHeight: '32px',
                flexShrink: 0,
                my: 1,
              }}
            >
              链接组
            </Box>
            {/* 使用 DragBrand 组件替换原有的品牌链接组 */}
            <DragBrand
              // @ts-expect-error 类型不匹配
              control={control}
              errors={errors}
              setIsEdit={setIsEdit}
            />
          </Box>
        </>
      )}
    </SettingCardItem>
  );
};

export default CardFooter;
