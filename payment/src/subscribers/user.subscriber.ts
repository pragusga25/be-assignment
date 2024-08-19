import { PrismaClient } from '@prisma/client';
import redisClient from '../redis';

const prisma = new PrismaClient();

export async function subscribeToUserCreation() {
  const subscriber = redisClient.duplicate();

  await subscriber.subscribe('user:created');

  subscriber.on('message', async (channel, message) => {
    console.log("Received message '%s' on channel '%s'", message, channel);
    if (channel === 'user:created') {
      const userData = JSON.parse(message);
      console.log('USER DATA', userData);
      await createUser(userData);
    }
  });
}

async function createUser(data: { id: string; email: string }) {
  try {
    await prisma.user.create({
      data,
    });
    console.log('User created:', data);
  } catch (error) {
    console.error('Error creating payment account:', error);
  }
}
