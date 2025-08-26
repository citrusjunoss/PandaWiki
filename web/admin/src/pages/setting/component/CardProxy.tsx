import { updateKnowledgeBase } from '@/api';
import { DomainKnowledgeBaseDetail } from '@/request/types';
import { SettingCardItem, FormItem } from './Common';

import {
  Box,
  FormControl,
  FormControlLabel,
  Radio,
  RadioGroup,
  Stack,
  TextField,
} from '@mui/material';
import { Message } from 'ct-mui';
import { useEffect, useState } from 'react';

const CardProxy = ({
  kb,
  refresh,
}: {
  kb: DomainKnowledgeBaseDetail;
  refresh: () => void;
}) => {
  const [isEdit, setIsEdit] = useState(false);
  const [hasProxy, setHasProxy] = useState(
    !!kb.access_settings?.trusted_proxies?.length,
  );
  const [proxyIP, setProxyIP] = useState(
    kb.access_settings?.trusted_proxies?.[0] || '',
  );

  const handleSave = () => {
    try {
      updateKnowledgeBase({
        id: kb.id,
        access_settings: {
          ...kb.access_settings,
          trusted_proxies: hasProxy ? [proxyIP] : null,
        },
      }).then(() => {
        Message.success('保存成功');
        setIsEdit(false);
        refresh();
      });
    } catch (e) {
      Message.error('保存失败');
    }
  };

  useEffect(() => {
    setHasProxy(!!kb.access_settings?.trusted_proxies?.length);
    setProxyIP(kb.access_settings?.trusted_proxies?.[0] || '');
  }, [kb]);

  return (
    <SettingCardItem
      title='前置反向代理'
      more={
        <Box
          sx={{
            flexGrow: 1,
            fontSize: 12,
            color: 'text.auxiliary',
            ml: 1,
            fontWeight: 'normal',
          }}
        >
          用于修正源 IP 获取错误的问题
        </Box>
      }
      isEdit={isEdit}
      onSubmit={handleSave}
    >
      <FormItem label='前置反向代理'>
        <FormControl>
          <RadioGroup
            value={hasProxy}
            onChange={e => {
              setHasProxy(e.target.value === 'true');
              if (proxyIP === '') {
                setProxyIP('0.0.0.0/0');
              }
              setIsEdit(true);
            }}
          >
            <Stack direction={'row'}>
              <FormControlLabel
                value={false}
                control={<Radio size='small' />}
                label='无前置反向代理'
              />
              <FormControlLabel
                value={true}
                control={<Radio size='small' />}
                label='有前置反向代理'
              />
            </Stack>
          </RadioGroup>
        </FormControl>
      </FormItem>

      {hasProxy && (
        <FormItem label='可信代理 IP'>
          <TextField
            fullWidth
            label='可信代理 IP'
            value={proxyIP}
            onChange={e => {
              setProxyIP(e.target.value);
              setIsEdit(true);
            }}
          />
        </FormItem>
      )}
    </SettingCardItem>
  );
};

export default CardProxy;
