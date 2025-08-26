import { DomainKnowledgeBaseDetail } from '@/request/types';
import { getApiProV1Prompt, postApiProV1Prompt } from '@/request/pro/Prompt';
import { Box, Slider, TextField } from '@mui/material';
import { useEffect, useMemo, useState } from 'react';
import { SettingCard, SettingCardItem, FormItem } from './Common';
import { Controller, useForm } from 'react-hook-form';
import { useAppSelector } from '@/store';
import { Message } from 'ct-mui';

interface CardAIProps {
  kb: DomainKnowledgeBaseDetail;
}

const CardAI = ({ kb }: CardAIProps) => {
  const [isEdit, setIsEdit] = useState(false);
  const { license } = useAppSelector(state => state.config);

  const { control, handleSubmit, setValue } = useForm({
    defaultValues: {
      interval: 0,
      content: '',
    },
  });

  const onSubmit = handleSubmit(async data => {
    await postApiProV1Prompt({
      kb_id: kb.id!,
      content: data.content,
    });

    Message.success('保存成功');
    setIsEdit(false);
  });

  const isPro = useMemo(() => {
    return license.edition === 1 || license.edition === 2;
  }, [license]);

  useEffect(() => {
    if (!kb.id || !isPro) return;
    getApiProV1Prompt({ kb_id: kb.id! }).then(res => {
      setValue('content', res.content || '');
    });
  }, [kb, isPro]);

  return (
    <SettingCard title='AI 设置'>
      <SettingCardItem title='智能问答' isEdit={isEdit} onSubmit={onSubmit}>
        <FormItem
          vertical
          tooltip={!isPro && '联创版和企业版可用'}
          extra={
            <Box
              sx={{
                fontSize: 12,
                color: 'primary.main',
                display: 'block',
                cursor: 'pointer',
              }}
              onClick={() => {
                setValue('content', '');
                setIsEdit(true);
              }}
            >
              重置为默认提示词
            </Box>
          }
          label='智能问答提示词'
        >
          <Controller
            control={control}
            name='content'
            render={({ field }) => (
              <TextField
                {...field}
                fullWidth
                disabled={!isPro}
                multiline
                rows={4}
                placeholder='智能问答提示词'
                onChange={e => {
                  field.onChange(e.target.value);
                  setIsEdit(true);
                }}
              />
            )}
          />
        </FormItem>

        <FormItem vertical label='连续提问时间间隔（敬请期待）'>
          <Controller
            control={control}
            name='interval'
            render={({ field }) => (
              <Slider
                {...field}
                disabled
                valueLabelDisplay='auto'
                min={200}
                max={300}
                step={5}
                sx={{
                  width: 432,
                  '& .MuiSlider-thumb': {
                    width: 16,
                    height: 16,
                    borderRadius: '50%',
                    backgroundColor: '#fff',
                    border: '2px solid currentColor',
                    '&:focus, &:hover, &.Mui-active, &.Mui-focusVisible': {
                      boxShadow: 'inherit',
                    },
                    '&::before': {
                      display: 'none',
                    },
                  },
                  '& .MuiSlider-track': {
                    bgcolor: 'primary.main',
                  },
                  '& .MuiSlider-rail': {
                    bgcolor: 'text.disabled',
                  },
                  '& .MuiSlider-valueLabel': {
                    lineHeight: 1.2,
                    fontSize: 12,
                    fontWeight: 'bold',
                    background: 'unset',
                    p: 0,
                    width: 24,
                    height: 24,
                    borderRadius: '50% 50% 50% 0',
                    bgcolor: 'primary.main',
                    transformOrigin: 'bottom left',
                    transform: 'translate(50%, -100%) rotate(-45deg) scale(0)',
                    '&::before': { display: 'none' },
                    '&.MuiSlider-valueLabelOpen': {
                      transform:
                        'translate(50%, -100%) rotate(-45deg) scale(1)',
                    },
                    '& > *': {
                      transform: 'rotate(45deg)',
                    },
                  },
                }}
                onChange={(e, value) => {
                  field.onChange(+value);
                  setIsEdit(true);
                }}
              />
            )}
          />
        </FormItem>
      </SettingCardItem>
    </SettingCard>
  );
};

export default CardAI;
