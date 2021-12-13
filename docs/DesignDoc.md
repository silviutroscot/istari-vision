# Design document for the Istari Vision project

## Architecture

- HTTP microservice / REST API
- TODO: Swagger file describing the API

### Alternatives

#### Publisher-subscriber infrastructure

##### Pros

- Loose coupling: the frontend and backend don't need to be able to know the DNS of eachother, they can use a queue
  where they post updates and the other consumes it; e.g. when we have a request the frontend creates a JSON
  representing the user's request, the back end consumes it and publishes the response to the same queue, on another
  topic and the front end can retrieve and render it.

##### Cons

- We have to maintain the queue(s)
- It adds more latency to communicate with the queue and then with the backend/frontend

## Testing

- We use a code coverage tool that will display the percentage of coverage we currently have
- We can add the coverage percentage in readme using codecov tool

### Unit testing

- For each strategy verify that the result is correct on various scenarios
- For the functions that validate if the user inputs are correct or not

### Integration testing

- Test for invalid inputs that the response code is 400 and the message can be displayed to the user so they can
  understand how to make a valid request
- Test for valid requests that the answer is calculated correctly and the response body, headers and status code are the
  expected ones
- Todo: add integration tests for the strategies

### Load testing

- Verify how many requests we can handle per second

### Performance testing

- Simulate multiple requests and compute the p50, p90, p95, p99 for the time to complete the requests

## Questions

- Should we run each strategy in a separate thread/go routine and merge the results in the end?
  - A: I don't think so, as the computation is very lightweight

## Deployment

- Use a VM on a cloud platform (e.g. AWS) to deploy the application (both backend and frontend)
- Automatically deploy on the VM, using CI pipeline (Github actions)

## Development process

- There is a `main` branch which contains the latest version of the code that was released in production and a `dev`
  branch which is the development branch
- Everyone who works on this project will create a branch from `dev` with the naming pattern `<username>/feature_name`
  e.g. `silviu/add-integration-tests-for-invalid-input`
- After finishing their task, they should make sure they rebased with the `dev` branch so they have the latest changes,
  before making a Pull Request(PR) to `dev`
- On PR the CI pipeline will start all the builds and run all tests, and if any step fails, the PR cannot be merged
- After the pipeline is successful, the PR needs at least one approval to be merged

## TODO in the future/if more time is provided

- Split the work left into epics and tasks and create a Jira/Trello board for the tasks, estimate their duration and
  prioritise them

- Create a Kubernetes cluster that will allow scaling the application and support more queries per second

- Have multiple environments for test, development, release candidate, preproduction, production and CI/CD pipelines and
  branches for them

- Have Ansible + Terraform scripts to configure and set up the environments automatically

- Use a queue for the requests, so we can handle them if they are too many to be handled

- If we add more complexity and the time to send a response increases we can send partial responses for a better user
  experience

- Have performance metrics so we can understand the bottlenecks of our system and how we can improve it + visualisation
  for them (using Grafana or any other tool) and aggregation (e.g. the median time to complete, p90, p95, p99, how many
  were finished because of time out, etc)

- Ingest and store the logs and have visualisaton and alerting for them; depending on how we want to work we can have a
  Slack integration, email integration, PagerDuty integration for alerting

- Have nightly performance tests that simulates a very large number of queries per second (QPS) to detect if we added
  any performance regression; to achieve this, we can have a Terraform file that creates and set up a VM and run the
  service + starts a client that sends progressively more requests to the service. After the test is finished our
  service writes the result to a file that can be uploaded to any BlobStorage system and also it compares the
  performance measurements with the baseline measurements and if we have regressions send an alert on email/slack to
  notify the team.

- Caching: store in a cache (e.g. Redis) the results for all investing strategies for a request and verify if we have
  already computed the return for that input (e.g. if we have a many users that want to see their potential return of
  investing 10 EGLD for 4 months and re-dele)

- A/B testing: use an A/B testing framework to determine for different implementations of the UI and new features how
  the customers are using it

- Consider using WebSockets to keep the TCP connection alive between frontend and backend (both ways) to enable faster
  updates of the live prices for MEX and EGLD we display in the front end
