import Redis from 'ioredis';
import { env } from './config';

const redisClient = new Redis(env.REDIS_URI);

export default redisClient;
