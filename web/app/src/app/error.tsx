'use client';
import ErrorPng from '@/assets/images/500.png';
import Footer from '@/components/footer';
import { Box, Stack } from '@mui/material';
import Image from 'next/image';

export default function Error({
  error,
  reset,
}: {
  error: Error & { digest?: string };
  reset: () => void;
}) {
  return (
    <Box
      sx={{
        position: 'relative',
        pt: 28,
        height: '100vh',
      }}
    >
      <Stack
        sx={{
          maxWidth: 1200,
          overflow: 'auto',
          pb: 6,
          mx: 'auto',
        }}
        justifyContent='center'
        alignItems='center'
      >
        <Image src={ErrorPng.src} alt='404' width={380} height={200} />
        <Stack
          gap={3}
          alignItems='center'
          sx={{ color: 'text.tertiary', fontSize: 14, mt: 3 }}
        >
          页面出错了 {error.digest}
        </Stack>
      </Stack>

      <Box
        sx={{
          height: 40,
          position: 'fixed',
          bottom: 0,
          left: 0,
          right: 0,
          zIndex: 1000,
        }}
      >
        <Footer showBrand={false} fullWidth={true} />
      </Box>
    </Box>
  );
}
