load 'create_docker_jobs'
describe '#create_job_manifest' do
  job = create_job_manifest('1.16', config, '1.16.1', 'job-1.16')

  it 'has a restart Policy of never' do
    expect(job[:spec][:template][:spec][:restartPolicy]).to eq 'Never'
  end
end
