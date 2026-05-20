const { Contract } = require('fabric-contract-api');

class CertificateContract extends Contract {

  async initLedger(ctx) {
    console.log('学业证书存证系统初始化成功！');
  }

  async uploadCertificate(ctx, certId, name, studentId, major, issueTime, hash) {
    const certificate = {
      certId,
      name,
      studentId,
      major,
      issueTime,
      hash,
      uploadTime: new Date().toLocaleString()
    };
    await ctx.stub.putState(certId, Buffer.from(JSON.stringify(certificate)));
    return `证书存证成功！ID：${certId}`;
  }

  async queryCertificate(ctx, certId) {
    const certificateBytes = await ctx.stub.getState(certId);
    if (!certificateBytes || certificateBytes.length === 0) {
      throw new Error(`证书ID ${certId} 不存在`);
    }
    return certificateBytes.toString();
  }

}

module.exports = CertificateContract;