import { KnowledgeBaseListItem, updateKnowledgeBase } from '@/api';
import { validateUrl } from '@/utils';
import { Box, Button, Stack, TextField } from '@mui/material';
import { Message } from 'ct-mui';
import { useEffect, useState } from 'react';
import CustomModal from '@/components/CustomModal';
import { SettingCardItem } from './Common';

const CardCustom = () => {
  const [customModalOpen, setCustomModalOpen] = useState(false);

  return (
    <SettingCardItem
      sx={{
        pb: '0 !important',
      }}
      title='前台网站样式个性化'
      extra={
        <Button
          size='small'
          variant='outlined'
          onClick={() => {
            setCustomModalOpen(true);
          }}
        >
          定制页面
        </Button>
      }
    >
      <CustomModal
        open={customModalOpen}
        onCancel={() => setCustomModalOpen(false)}
      />
    </SettingCardItem>
  );
};

export default CardCustom;
