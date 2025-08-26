import { DomainKnowledgeBaseDetail } from '@/request/types';
import {
  Box,
  Button,
  FormControl,
  FormControlLabel,
  Link,
  Radio,
  RadioGroup,
  Stack,
  TextField,
} from '@mui/material';
import { useEffect, useState } from 'react';
import { FormItem, SettingCardItem } from './Common';

const CardRobotApi = ({ kb }: { kb: DomainKnowledgeBaseDetail }) => {
  const [isEdit, setIsEdit] = useState(false);
  const [isEnabled, setIsEnabled] = useState(false);

  useEffect(() => {}, [kb]);

  return (
    <SettingCardItem
      title='问答机器人 API（敬请期待）'
      isEdit={isEdit}
      more={
        <Link
          component='a'
          href='https://pandawiki.docs.baizhi.cloud/node/01971b60-100e-7b23-9385-e36763df5c0a'
          target='_blank'
          sx={{
            fontSize: 14,
            ml: 1,
            textDecoration: 'none',
            fontWeight: 'normal',
            '&:hover': {
              fontWeight: 'bold',
            },
          }}
        >
          使用方法
        </Link>
      }
      onSubmit={() => {}}
    >
      <FormItem label='问答机器人 API'>
        <FormControl>
          <RadioGroup
            value={isEnabled}
            onChange={e => {
              setIsEnabled(e.target.value === 'true');
              setIsEdit(e.target.value === 'true');
            }}
          >
            <Stack direction={'row'}>
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
            </Stack>
          </RadioGroup>
        </FormControl>
      </FormItem>

      {isEnabled && (
        <>
          <FormItem label='API Token'>
            <TextField fullWidth label='API Token' placeholder={'API Token'} />
          </FormItem>
          <FormItem label='API 调用地址'>
            <TextField
              fullWidth
              label='API 调用地址'
              placeholder={'API 调用地址'}
            />
          </FormItem>
        </>
      )}
    </SettingCardItem>
  );
};

export default CardRobotApi;
