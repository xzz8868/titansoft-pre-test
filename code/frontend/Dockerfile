# Use nginx 1.27 as the base image
FROM nginx:1.27

# Copy the local nginx configuration file to the container
COPY nginx.conf /etc/nginx/nginx.conf

# Copy the local static web files to the appropriate location in the container
COPY html /usr/share/nginx/html

# Command to run when the container starts up
CMD ["nginx", "-g", "daemon off;"]