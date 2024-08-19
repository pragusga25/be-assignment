import fastifyCookie from '@fastify/cookie';
import type { TypeBoxTypeProvider } from '@fastify/type-provider-typebox';
import Fastify from 'fastify';
import { verifyToken } from './middlewares/auth.middleware';
import paymentRoutes from './routes/payment.route';
import { subscribeToUserCreation } from './subscribers/user.subscriber';

const server = Fastify({
  logger: true,
}).withTypeProvider<TypeBoxTypeProvider>();

server.register(fastifyCookie);

server.addHook('preHandler', verifyToken);

server.register(paymentRoutes, { prefix: '/api/payments' });

const start = async () => {
  try {
    await subscribeToUserCreation();
    await server.listen({ port: 3000, host: '0.0.0.0' });
    console.log('Server is running on http://localhost:3000');
  } catch (err) {
    server.log.error(err);
    process.exit(1);
  }
};

start();
