import CardRobotApi from './CardRobotApi';
import CardRobotDing from './CardRobotDing';
import CardRobotDiscord from './CardRobotDiscord';
import CardRobotFeishu from './CardRobotFeishu';
import CardRobotWebComponent from './CardRobotWebComponent';
import CardRobotWechatOfficeAccount from './CardRobotWechatOfficeAccount';
import CardRobotWecom from './CardRobotWecom';
import CardRobotWecomService from './CardRobotWecomService';
import { DomainKnowledgeBaseDetail } from '@/request/types';
import { SettingCard } from './Common';

const CardRobot = ({
  kb,
  url,
}: {
  kb: DomainKnowledgeBaseDetail;
  url: string;
}) => {
  return (
    <SettingCard title='问答机器人'>
      <CardRobotWebComponent kb={kb} />
      <CardRobotApi kb={kb} />
      <CardRobotDing kb={kb} />
      <CardRobotFeishu kb={kb} />
      <CardRobotWechatOfficeAccount kb={kb} url={url} />
      <CardRobotWecom kb={kb} url={url} />
      <CardRobotWecomService kb={kb} url={url} />
      <CardRobotDiscord kb={kb} />
    </SettingCard>
  );
};

export default CardRobot;
