---
description: Installing lakeFS is easy. This section covers how to install lakeFS using docker compose
---

# Install lakeFS

{% hint style="info" %}
**Note:** The quickstart section is for learning purposes. The installations below will not persist your data.
Instead, it will spin-up a database in a docker container, which will be discarded later.

For a production suitable deployment, see [Deploy](/deploy).
{% endhint %}

## Using docker-compose

Other quickstart methods can be found [here](more_quickstart_options.md).

To run a local lakeFS instance using [Docker Compose](https://docs.docker.com/compose/):

1. Ensure you have Docker & Docker Compose installed on your computer, and that compose version is 1.25.04 or higher. For more information, please see this [issue](https://github.com/treeverse/lakeFS/issues/894). 
1. Run the following command in your terminal:

   ```bash
   curl https://compose.lakefs.io | docker-compose -f - up
   ```

1. Check your installation by opening [http://127.0.0.1:8000/setup](http://127.0.0.1:8000/setup) in your web browser.

### Next steps

Now that your lakeFS is running, try [creating a repository](repository.md).
