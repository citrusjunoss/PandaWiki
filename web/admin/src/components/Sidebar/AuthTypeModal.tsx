import HelpCenter from '@/assets/json/help-center.json';
import IconUpgrade from '@/assets/json/upgrade.json';
import {
  Box,
  Button,
  Stack,
} from '@mui/material';
import { Modal } from 'ct-mui';
import LottieIcon from '../LottieIcon';

interface AuthTypeModalProps {
  open: boolean;
  onClose: () => void;
  curVersion: string;
  latestVersion: string;
}

const AuthTypeModal = ({
  open,
  onClose,
  curVersion,
  latestVersion,
}: AuthTypeModalProps) => {

  return (
    <>
      <Modal
        open={open}
        footer={null}
        title='关于 PandaWiki'
        onCancel={onClose}
      >
        <Stack gap={1} sx={{ fontSize: 14, lineHeight: '32px' }}>
          <Stack direction={'row'} alignItems={'center'}>
            <Box sx={{ width: 120, flexShrink: 0 }}>当前版本</Box>
            <Stack direction={'row'} alignItems={'center'} gap={2}>
              <Box sx={{ fontFamily: 'GBold', minWidth: 50 }}>{curVersion}</Box>
              {latestVersion === `v${curVersion}` ? (
                <Box sx={{ color: 'text.auxiliary', fontSize: 12 }}>
                  已是最新版本，无需更新
                </Box>
              ) : (
                <Button
                  size='small'
                  startIcon={
                    <Box>
                      <LottieIcon
                        id='version'
                        src={latestVersion === '' ? HelpCenter : IconUpgrade}
                        style={{ width: 16, height: 16, display: 'flex' }}
                      />
                    </Box>
                  }
                  onClick={() => {
                    window.open(
                      'https://pandawiki.docs.baizhi.cloud/node/01971615-05b8-7924-9af7-15f73784f893',
                    );
                  }}
                >
                  立即更新
                </Button>
              )}
            </Stack>
          </Stack>
          <Stack direction={'row'} alignItems={'center'}>
            <Box sx={{ width: 120, flexShrink: 0 }}>产品型号</Box>
            <Stack direction={'row'} alignItems={'center'} gap={2}>
              <Box sx={{ minWidth: 50 }}>开源版</Box>
            </Stack>
          </Stack>
        </Stack>
      </Modal>
    </>
  );
};

export default AuthTypeModal;
