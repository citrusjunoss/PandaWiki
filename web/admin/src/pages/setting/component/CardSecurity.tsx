import React, { useMemo, useState, useEffect } from 'react';
import { SettingCardItem, SettingCard, FormItem } from './Common';
import { useAppSelector } from '@/store';
import { useForm, Controller } from 'react-hook-form';
import { putApiV1App } from '@/request/App';
import {
  Radio,
  RadioGroup,
  FormControlLabel,
  TextField,
  styled,
  Autocomplete,
} from '@mui/material';
import { Message } from 'ct-mui';
import {
  DomainAppDetailResp,
  DomainKnowledgeBaseDetail,
} from '@/request/types';
import { getApiProV1Block, postApiProV1Block } from '@/request/pro/Block';

const StyledRadioLabel = styled('div')(({ theme }) => ({
  width: 100,
}));

const WatermarkForm = ({ data }: { data?: DomainAppDetailResp }) => {
  const { license } = useAppSelector(state => state.config);
  const [watermarkIsEdit, setWatermarkIsEdit] = useState(false);
  const {
    control,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm({
    defaultValues: {
      watermark_enable: data?.settings?.watermark_enable ?? null,
      watermark_content: data?.settings?.watermark_content ?? '',
    },
  });

  const watermarkEnable = watch('watermark_enable');
  const isEnterprise = useMemo(() => {
    return license.edition === 2;
  }, [license]);

  const handleSaveWatermark = handleSubmit(values => {
    if (!data?.id || values.watermark_enable === null) return;
    putApiV1App(
      { id: data.id },
      {
        settings: {
          ...data?.settings,
          watermark_enable: values.watermark_enable,
          watermark_content: values.watermark_content,
        },
      },
    ).then(() => {
      Message.success('保存成功');
      setWatermarkIsEdit(false);
    });
  });

  useEffect(() => {
    if (!data) return;
    setValue('watermark_enable', data.settings?.watermark_enable ?? null);
    setValue('watermark_content', data.settings?.watermark_content ?? '');
  }, [data]);

  return (
    <SettingCardItem
      title='水印'
      isEdit={watermarkIsEdit}
      onSubmit={handleSaveWatermark}
    >
      <FormItem label='水印开关' tooltip={!isEnterprise && '企业版可用'}>
        <Controller
          control={control}
          name='watermark_enable'
          render={({ field }) => (
            <RadioGroup
              row
              {...field}
              onChange={e => {
                setWatermarkIsEdit(true);
                field.onChange(e.target.value === 'true');
              }}
            >
              <FormControlLabel
                value={true}
                control={<Radio size='small' disabled={!isEnterprise} />}
                label={<StyledRadioLabel>启用</StyledRadioLabel>}
              />
              <FormControlLabel
                value={false}
                control={<Radio size='small' disabled={!isEnterprise} />}
                label={<StyledRadioLabel>禁用</StyledRadioLabel>}
              />
            </RadioGroup>
          )}
        />
      </FormItem>
      {watermarkEnable && (
        <FormItem label='水印内容' sx={{ alignItems: 'flex-start' }}>
          <Controller
            control={control}
            name='watermark_content'
            render={({ field }) => (
              <TextField
                fullWidth
                {...field}
                placeholder='请输入水印内容, 支持多行输入'
                multiline
                minRows={2}
                disabled={!isEnterprise}
                onChange={e => {
                  setWatermarkIsEdit(true);
                  field.onChange(e.target.value);
                }}
              />
            )}
          />
        </FormItem>
      )}
    </SettingCardItem>
  );
};
const VerificationForm = ({ data }: { data?: DomainAppDetailResp }) => {
  const { license } = useAppSelector(state => state.config);
  const [isEdit, setIsEdit] = useState(false);
  const {
    control,
    handleSubmit,
    setValue,
    watch,
    formState: { errors },
  } = useForm({
    defaultValues: {
      verification_enable: null,
      comment_verification_enable: null,
    },
  });

  const verificationEnable = watch('verification_enable');
  const isEnterprise = useMemo(() => {
    return license.edition === 2;
  }, [license]);

  const handleSaveVerification = handleSubmit(values => {
    // if (!data?.id || values.verification_enable === null) return;
    // putApiV1App(
    //   { id: data.id },
    //   {
    //     settings: {
    //       ...data?.settings,
    //       verification_enable: values.verification_enable,
    //       verification_content: values.verification_content,
    //     },
    //   },
    // ).then(() => {
    //   Message.success('保存成功');
    //   setIsEdit(false);
    // });
  });

  useEffect(() => {
    if (!data) return;
    // setValue('verification_enable', data.settings?.verification_enable ?? null);
    // setValue('verification_content', data.settings?.verification_content ?? '');
  }, [data]);

  return (
    <SettingCardItem
      title='人机验证（敬请期待）'
      isEdit={isEdit}
      onSubmit={handleSaveVerification}
    >
      <FormItem label='问答'>
        <Controller
          control={control}
          name='verification_enable'
          render={({ field }) => (
            <RadioGroup
              row
              {...field}
              onChange={e => {
                setIsEdit(true);
                field.onChange(e.target.value === 'true');
              }}
            >
              <FormControlLabel
                value={true}
                control={<Radio size='small' disabled />}
                label={<StyledRadioLabel>需要验证码</StyledRadioLabel>}
              />
              <FormControlLabel
                value={false}
                control={<Radio size='small' disabled />}
                label={<StyledRadioLabel>无需验证码</StyledRadioLabel>}
              />
            </RadioGroup>
          )}
        />
      </FormItem>
      <FormItem label='评论'>
        <Controller
          control={control}
          name='comment_verification_enable'
          render={({ field }) => (
            <RadioGroup
              row
              {...field}
              onChange={e => {
                setIsEdit(true);
                field.onChange(e.target.value === 'true');
              }}
            >
              <FormControlLabel
                value={true}
                control={<Radio size='small' disabled />}
                label={<StyledRadioLabel>需要验证码</StyledRadioLabel>}
              />
              <FormControlLabel
                value={false}
                control={<Radio size='small' disabled />}
                label={<StyledRadioLabel>无需验证码</StyledRadioLabel>}
              />
            </RadioGroup>
          )}
        />
      </FormItem>
    </SettingCardItem>
  );
};

const KeywordsForm = ({ kb }: { kb: DomainKnowledgeBaseDetail }) => {
  const { license } = useAppSelector(state => state.config);
  const [questionInputValue, setQuestionInputValue] = useState('');
  const [isEdit, setIsEdit] = useState(false);
  const isEnterprise = useMemo(() => {
    return license.edition === 2;
  }, [license]);

  const { control, handleSubmit, setValue } = useForm({
    defaultValues: {
      interval: 0,
      content: '',
      block_words: [] as string[],
    },
  });

  const onSubmit = handleSubmit(async data => {
    await postApiProV1Block({
      kb_id: kb.id!,
      block_words: data.block_words,
    });

    Message.success('保存成功');
    setIsEdit(false);
  });

  useEffect(() => {
    if (!kb.id || !isEnterprise) return;
    getApiProV1Block({ kb_id: kb.id! }).then(res => {
      setValue('block_words', res.words || []);
    });
  }, [kb, isEnterprise]);

  return (
    <SettingCardItem title='内容合规' isEdit={isEdit} onSubmit={onSubmit}>
      <FormItem
        vertical
        tooltip={!isEnterprise && '企业版可用'}
        label='屏蔽 AI 问答中的关键字'
      >
        <Controller
          control={control}
          name='block_words'
          render={({ field }) => (
            <Autocomplete
              {...field}
              multiple
              freeSolo
              inputValue={questionInputValue}
              options={[]}
              fullWidth
              disabled={!isEnterprise}
              onInputChange={(_, value) => {
                setQuestionInputValue(value);
              }}
              onChange={(_, newValue) => {
                setIsEdit(true);

                const newValues = [...new Set(newValue as string[])];
                field.onChange(newValues);
              }}
              renderInput={params => (
                <TextField
                  {...params}
                  placeholder='屏蔽 AI 问答中的关键字，可输入多个，回车确认'
                  variant='outlined'
                  onBlur={() => {
                    // 失去焦点时自动添加当前输入的值
                    const trimmedValue = questionInputValue.trim();
                    if (trimmedValue && !field.value.includes(trimmedValue)) {
                      setIsEdit(true);
                      field.onChange([...field.value, trimmedValue]);
                      setQuestionInputValue('');
                    }
                  }}
                />
              )}
            />
          )}
        />
      </FormItem>
    </SettingCardItem>
  );
};

const CardSecurity = ({
  data,
  kb,
}: {
  data?: DomainAppDetailResp;
  kb: DomainKnowledgeBaseDetail;
}) => {
  return (
    <SettingCard title='安全'>
      <VerificationForm data={data} />
      <WatermarkForm data={data} />
      <KeywordsForm kb={kb} />
    </SettingCard>
  );
};

export default CardSecurity;
