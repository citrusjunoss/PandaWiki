import { updateAppDetail, WelcomeSetting } from '@/api';
import { getApiV1NodeRecommendNodes } from '@/request/Node';
import {
  DomainAppDetailResp,
  DomainRecommendNodeListResp,
} from '@/request/types';
import DragRecommend from '@/components/Drag/DragRecommend';
import { FreeSoloAutocomplete } from '@/components/FreeSoloAutocomplete';
import { useCommitPendingInput } from '@/hooks';
import { useAppSelector } from '@/store';
import { Box, Button, Stack, TextField } from '@mui/material';
import { Icon, Message } from 'ct-mui';
import { useEffect, useState } from 'react';
import { Controller, useForm } from 'react-hook-form';
import AddRecommendContent from './AddRecommendContent';
import { FormItem, SettingCardItem } from './Common';

interface CardWebWelcomeProps {
  id: string;
  data: DomainAppDetailResp;
  refresh: (value: WelcomeSetting) => void;
}

const CardWebWelcome = ({ id, data, refresh }: CardWebWelcomeProps) => {
  const { kb_id } = useAppSelector(state => state.config);
  const [sorted, setSorted] = useState<DomainRecommendNodeListResp[]>([]);
  const [isEdit, setIsEdit] = useState(false);
  const [open, setOpen] = useState(false);
  const {
    control,
    handleSubmit,
    formState: { errors },
    watch,
    setValue,
  } = useForm<WelcomeSetting>({
    defaultValues: {
      welcome_str: '',
      search_placeholder: '',
      recommend_questions: [],
      recommend_node_ids: [],
    },
  });

  const recommend_questions = watch('recommend_questions') || [];
  const recommend_node_ids = watch('recommend_node_ids') || [];

  const recommendQuestionsField = useCommitPendingInput<string>({
    value: recommend_questions,
    setValue: value => {
      setIsEdit(true);
      setValue('recommend_questions', value);
    },
  });

  const onSubmit = handleSubmit(value => {
    // @ts-expect-error 类型不匹配
    updateAppDetail({ id }, { settings: { ...data.settings, ...value } }).then(
      () => {
        refresh(value);
        Message.success('保存成功');
        setIsEdit(false);
      },
    );
  });

  const nodeRec = () => {
    getApiV1NodeRecommendNodes({ kb_id, node_ids: recommend_node_ids }).then(
      res => {
        setSorted(res);
      },
    );
  };

  useEffect(() => {
    if (recommend_node_ids.length > 0) {
      nodeRec();
    } else {
      setSorted([]);
    }
  }, [recommend_node_ids]);

  useEffect(() => {
    setSorted(data.recommend_nodes || []);
    setValue('welcome_str', data.settings?.welcome_str || '');
    setValue('search_placeholder', data.settings?.search_placeholder || '');
    setValue('recommend_questions', data.settings?.recommend_questions || []);
    setValue('recommend_node_ids', data.settings?.recommend_node_ids || []);
  }, [data]);

  return (
    <SettingCardItem title='欢迎页面' isEdit={isEdit} onSubmit={onSubmit}>
      <FormItem label='欢迎标语'>
        {' '}
        <Controller
          control={control}
          name='welcome_str'
          render={({ field }) => (
            <TextField
              fullWidth
              {...field}
              placeholder='输入欢迎标语'
              error={!!errors.welcome_str}
              helperText={errors.welcome_str?.message}
              onChange={event => {
                setIsEdit(true);
                field.onChange(event);
              }}
            />
          )}
        />
      </FormItem>
      <FormItem label='搜索框提示文字'>
        <Controller
          control={control}
          name='search_placeholder'
          render={({ field }) => (
            <TextField
              fullWidth
              {...field}
              placeholder='输入搜索框提示文字'
              error={!!errors.search_placeholder}
              helperText={errors.search_placeholder?.message}
              onChange={event => {
                setIsEdit(true);
                field.onChange(event);
              }}
            />
          )}
        />
      </FormItem>
      <FormItem label='推荐问题'>
        <FreeSoloAutocomplete
          placeholder='回车确认，填写下一个推荐问题'
          {...recommendQuestionsField}
        />
      </FormItem>
      <FormItem label='推荐内容' vertical>
        <Box sx={{ width: '100%' }}>
          <Box sx={{ mb: 1 }}>
            <DragRecommend
              data={sorted || []}
              refresh={nodeRec}
              onChange={value => {
                setIsEdit(true);
                setValue(
                  'recommend_node_ids',
                  value.map(item => item.id!),
                );
              }}
            />
          </Box>
          <Button
            size='small'
            onClick={() => setOpen(true)}
            startIcon={
              <Icon type='icon-add' sx={{ fontSize: '12px !important' }} />
            }
          >
            添加卡片
          </Button>
        </Box>
      </FormItem>

      <AddRecommendContent
        open={open}
        selected={recommend_node_ids}
        onChange={(value: string[]) => {
          setIsEdit(true);
          setValue('recommend_node_ids', value);
        }}
        onClose={() => setOpen(false)}
      />
    </SettingCardItem>
  );
};
export default CardWebWelcome;
