import { FooterProvider } from '@/components/footer';
import Header from '@/components/header';
import Home from '@/views/home';
import Catalog from '@/views/node/Catalog';
import { Box } from '@mui/material';
import WaterMarkProvider from '@/components/watermark/WaterMarkProvider';

const HomePage = () => {
  return (
    <WaterMarkProvider>
      <Box
        sx={{
          position: 'relative',
          bgcolor: 'background.default',
        }}
      >
        <Catalog />
        <Header />
        <Home />
        <FooterProvider />
      </Box>
    </WaterMarkProvider>
  );
};

export default HomePage;
