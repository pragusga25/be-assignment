import { PaymentAccountType } from '@prisma/client';
import { type Static, Type } from '@sinclair/typebox';
import type { FastifyInstance } from 'fastify';
import * as PaymentService from '../services/payment.service';

export default async function (fastify: FastifyInstance) {
  const DepositSchema = Type.Object({
    paymentAccountId: Type.String(),
    amount: Type.Number(),
  });

  fastify.post<{ Body: Static<typeof DepositSchema> }>('/deposit', {
    schema: {
      body: DepositSchema,
    },
    handler: async (request, reply) => {
      const result = await PaymentService.deposit({
        ...request.body,
        // biome-ignore lint/style/noNonNullAssertion: <explanation>
        userId: request!.user!.id,
      });
      reply.send(result);
    },
  });

  const WithdrawSchema = Type.Object({
    paymentAccountId: Type.String(),
    amount: Type.Number(),
  });
  fastify.post<{
    Body: Static<typeof WithdrawSchema>;
  }>('/withdraw', {
    schema: {
      body: WithdrawSchema,
    },
    handler: async (request, reply) => {
      const result = await PaymentService.withdraw({
        ...request.body,
        // biome-ignore lint/style/noNonNullAssertion: <explanation>
        userId: request.user!.id,
      });

      reply.send(result);
    },
  });

  const SendSchema = Type.Object({
    fromPaymentAccountId: Type.String(),
    toPaymentAccountId: Type.String(),
    amount: Type.Number(),
  });
  fastify.post<{
    Body: Static<typeof SendSchema>;
  }>('/send', {
    schema: {
      body: SendSchema,
    },
    handler: async (request, reply) => {
      const result = await PaymentService.send(
        // biome-ignore lint/style/noNonNullAssertion: <explanation>
        { ...request.body, userId: request!.user!.id }
      );
      reply.send(result);
    },
  });

  fastify.get('/histories', async (request, reply) => {
    const histories = await PaymentService.getPaymentHistories(
      // biome-ignore lint/style/noNonNullAssertion: <explanation>
      request!.user!.id
    );
    reply.send(histories);
  });

  const CreatePaymentAccountSchema = Type.Object({
    paymentAccountId: Type.String(),
    paymentAccountType: Type.Enum(PaymentAccountType),
    balance: Type.Number(),
  });

  fastify.post<{ Body: Static<typeof CreatePaymentAccountSchema> }>(
    '/accounts',
    {
      schema: {
        body: CreatePaymentAccountSchema,
      },
      handler: async (request, reply) => {
        const result = await PaymentService.createPaymentAccount({
          ...request.body,
          // biome-ignore lint/style/noNonNullAssertion: <explanation>
          userId: request!.user!.id,
        });
        reply.send(result);
      },
    }
  );

  fastify.get('/accounts', async (request, reply) => {
    const histories = await PaymentService.getMyPaymentAccounts(
      // biome-ignore lint/style/noNonNullAssertion: <explanation>
      request!.user!.id
    );
    reply.send(histories);
  });
}
