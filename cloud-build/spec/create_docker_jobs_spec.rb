describe '#create_job_manifest' do
  let(:job) do
    load 'create_docker_jobs'
    create_job_manifest('1.16', config, '1.16.1', 'job-1.16')
  end

  it 'has a restart Policy of never' do
    expect(job.spec.template.spec.restart_policy).to eq 'Never'
  end
end
