import { PaymentStatus, PrismaClient } from '@prisma/client';
import type {
  ICreatePaymentAccount,
  IDepositData,
  ISendData,
} from '../interfaces';

const prisma = new PrismaClient();

export async function deposit(data: IDepositData) {
  const { amount, paymentAccountId, userId } = data;
  return prisma.$transaction(async (tx) => {
    const account = await tx.paymentAccount.findFirst({
      where: { paymentAccountId, userId },
    });
    if (!account) throw new Error('Account not found');

    await tx.paymentAccount.update({
      where: { paymentAccountId },
      data: { balance: account.balance + amount },
    });

    return tx.paymentHistory.create({
      data: {
        toPaymentAccountId: paymentAccountId,
        fromPaymentAccountId: paymentAccountId,
        amount,
        status: PaymentStatus.COMPLETED,
      },
    });
  });
}

export async function withdraw(data: IDepositData) {
  const { amount, paymentAccountId, userId } = data;
  return prisma.$transaction(async (tx) => {
    const account = await tx.paymentAccount.findUnique({
      where: { paymentAccountId, userId },
    });
    if (!account) throw new Error('Account not found');
    if (account.balance < amount) throw new Error('Insufficient funds');

    await tx.paymentAccount.update({
      where: { paymentAccountId },
      data: { balance: account.balance - amount },
    });

    return tx.paymentHistory.create({
      data: {
        toPaymentAccountId: paymentAccountId,
        fromPaymentAccountId: paymentAccountId,
        amount: -amount,
        status: PaymentStatus.COMPLETED,
      },
    });
  });
}

export async function send(data: ISendData) {
  const { amount, fromPaymentAccountId, toPaymentAccountId, userId } = data;
  return prisma.$transaction(async (tx) => {
    const fromAccount = await tx.paymentAccount.findUnique({
      where: { paymentAccountId: fromPaymentAccountId, userId },
    });
    const toAccount = await tx.paymentAccount.findUnique({
      where: { paymentAccountId: toPaymentAccountId },
    });

    if (!fromAccount || !toAccount)
      throw new Error('One or both accounts not found');
    if (fromAccount.balance < amount) throw new Error('Insufficient funds');

    await tx.paymentAccount.update({
      where: { paymentAccountId: fromPaymentAccountId },
      data: { balance: fromAccount.balance - amount },
    });

    await tx.paymentAccount.update({
      where: { paymentAccountId: toPaymentAccountId },
      data: { balance: toAccount.balance + amount },
    });

    return tx.paymentHistory.create({
      data: {
        fromPaymentAccountId: fromPaymentAccountId,
        toPaymentAccountId: toPaymentAccountId,
        amount,
        status: PaymentStatus.COMPLETED,
      },
    });
  });
}

export async function getPaymentHistories(userId: string) {
  return prisma.paymentHistory.findMany({
    where: {
      OR: [
        {
          fromAccount: {
            userId,
          },
        },
        { toAccount: { userId } },
      ],
    },

    select: {
      amount: true,
      status: true,
      createdAt: true,
      fromAccount: {
        select: {
          paymentAccountId: true,
          user: {
            select: {
              email: true,
            },
          },
        },
      },
      toAccount: {
        select: {
          paymentAccountId: true,
          user: {
            select: {
              email: true,
            },
          },
        },
      },
    },
  });
}

export async function getMyPaymentAccounts(userId: string) {
  return prisma.paymentAccount.findMany({
    where: { userId },
  });
}

export async function createPaymentAccount(data: ICreatePaymentAccount) {
  return prisma.paymentAccount.create({
    data,
  });
}
