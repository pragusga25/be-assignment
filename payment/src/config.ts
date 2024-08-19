import dotenv from 'dotenv';

dotenv.config();

export const env = {
  DATABASE_URL: process.env.DATABASE_URL as string,
  REDIS_URI: process.env.REDIS_URI as string,
  JWKS_URI: process.env.JWKS_URI as string,
  PORT: Number.parseInt(process.env.PORT || '3000', 10),
};
