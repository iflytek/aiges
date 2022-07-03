class UnknownCudaRCDistro(Exception):
    """An exception that is raised when a match connot be made against a distro from Shipit global.json"""

    pass


class RequestsRetry(Exception):
    """An exception to handle retries for requests http gets"""

    pass


class ImageRegistryLoginRetry(Exception):
    """An exception to handle retries for container registry login"""

    pass


class ImagePushRetry(Exception):
    """An exception to handle retries for pushing container images"""

    pass


class ImageDeleteRetry(Exception):
    """An exception to handle retries for image deletion"""

    pass

