// Next.js API route support: https://nextjs.org/docs/api-routes/introduction
import type { NextApiRequest, NextApiResponse } from 'next';
import { Gage } from '../../types';
import { UserConfig } from '../../api/database/database';
import { handleRequest } from '../../lib/handleRequest';

type Data = {
  gages: Gage[] | undefined;
};

export default async function handler(
  req: NextApiRequest,
  res: NextApiResponse<Data>,
) {
  await handleRequest(UserConfig, req, res);
}