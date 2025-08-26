import { updateAppDetail } from '@/api';
import { ThemeAndStyleSetting, ThemeMode } from '@/api/type';
import UploadFile from '@/components/UploadFile';
import { MenuItem, Select } from '@mui/material';
import { Message } from 'ct-mui';
import { useEffect, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import { FormItem, SettingCardItem } from './Common';
import { DomainAppDetailResp } from '@/request/types';

interface CardStyleProps {
  id: string;
  data: DomainAppDetailResp;
  refresh: (value: ThemeMode & ThemeAndStyleSetting) => void;
}

const CardStyle = ({ id, data, refresh }: CardStyleProps) => {
  const [isEdit, setIsEdit] = useState(false);
  const { control, handleSubmit, setValue } = useForm<
    ThemeMode & ThemeAndStyleSetting
  >({
    defaultValues: {
      theme_mode: 'light',
      bg_image: '',
    },
  });

  const onSubmit = (value: ThemeMode & ThemeAndStyleSetting) => {
    updateAppDetail(
      { id },
      {
        // @ts-expect-error 类型不匹配
        settings: {
          ...data.settings,
          theme_mode: value.theme_mode,
          theme_and_style: {
            ...data.settings?.theme_and_style,
            bg_image: value.bg_image,
          },
        },
      },
    ).then(() => {
      refresh(value);
      Message.success('保存成功');
      setIsEdit(false);
    });
  };

  useEffect(() => {
    setValue('theme_mode', data.settings?.theme_mode as 'light' | 'dark');
    setValue('bg_image', data.settings?.theme_and_style?.bg_image || '');
  }, [data]);

  return (
    <SettingCardItem
      title='样式与风格'
      isEdit={isEdit}
      onSubmit={handleSubmit(onSubmit)}
    >
      <FormItem label='配色方案'>
        <Controller
          control={control}
          name='theme_mode'
          render={({ field }) => (
            <Select
              {...field}
              sx={{ width: '100%', height: 52 }}
              onChange={e => {
                field.onChange(e.target.value as 'light' | 'dark');
                setIsEdit(true);
              }}
            >
              <MenuItem value='light'>浅色模式</MenuItem>
              <MenuItem value='dark'>深色模式</MenuItem>
            </Select>
          )}
        />
      </FormItem>

      <FormItem label='背景图片'>
        <Controller
          control={control}
          name='bg_image'
          render={({ field }) => (
            <UploadFile
              {...field}
              id='bg_image'
              type='url'
              accept='image/*'
              width={80}
              onChange={url => {
                field.onChange(url);
                setIsEdit(true);
              }}
            />
          )}
        />{' '}
      </FormItem>
    </SettingCardItem>
  );
};

export default CardStyle;
