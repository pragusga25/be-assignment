import type { FastifyReply, FastifyRequest } from 'fastify';
import JsonWebToken, {
  type JwtHeader,
  type SigningKeyCallback,
  type JwtPayload,
} from 'jsonwebtoken';
import jwksClient from 'jwks-rsa';
import { env } from '../config';

const client = jwksClient({
  jwksUri: env.JWKS_URI,
});

function getKey(header: JwtHeader, callback: SigningKeyCallback) {
  client.getSigningKey(header.kid, (err, key) => {
    const signingKey = key?.getPublicKey();
    callback(err, signingKey);
  });
}

export async function verifyToken(
  request: FastifyRequest,
  reply: FastifyReply
) {
  const token = request.cookies.access_token;

  if (!token) {
    reply.code(401).send({ error: 'No token provided' });
    return;
  }

  try {
    const decoded = await new Promise<JwtPayload>((resolve, reject) => {
      JsonWebToken.verify(token, getKey, {}, (err, decoded) => {
        if (err) reject(err);
        else resolve(decoded as JwtPayload);
      });
    });

    request.user = {
      id: decoded.sub as string,
    };
  } catch (error) {
    reply.code(401).send({ error: 'Invalid token' });
  }
}
